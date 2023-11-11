package models

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"sync"
	"v/pkg/config"
	protocol "v/protocol/go_protocol"

	stt "cloud.google.com/go/speech/apiv1"
	sttpb "cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pion/webrtc/v3"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type TranscriberModel struct {
	db *gorm.DB
	ctx context.Context
	client *openai.Client
	cancel context.CancelFunc

	oggWriter     *io.PipeWriter
	oggReader     *io.PipeReader
	oggSerializer *oggwriter.OggWriter
	lock sync.Mutex
	closeCh chan struct{}
	results chan TranscriptionResult
	language Language
	speechClient *stt.Client
	rtpCodec webrtc.RTPCodecParameters
}

type TranscriptionResult struct {
	Error   error
	Text    string
	IsFinal bool
}

func NewTranscriberModel(language protocol.Language,rtpCodec webrtc.RTPCodecParameters, speechClient *stt.Client) (*TranscriberModel, error) {
	//token := config.Conf.Openai.Token
	//client := openai.NewClient(token)
	ctx :=context.Background()
		if !strings.EqualFold(rtpCodec.MimeType, "audio/opus") {
		return nil, errors.New("only opus is supported")
	}

	oggReader, oggWriter := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	t := &TranscriberModel{
		db: config.App.DB,
		//client: client,
		ctx: ctx,

		cancel:   cancel,
		rtpCodec: rtpCodec,
		//sb:           samplebuilder.New(200, &codecs.OpusPacket{}, rtpCodec.ClockRate),
		oggReader:    oggReader,
		oggWriter:    oggWriter,
		language:     NewLanguageModel(language),
		speechClient: speechClient,
		results:      make(chan TranscriptionResult),
		closeCh:      make(chan struct{}),
	}
	go t.start()
	return t, nil

}


func (t *TranscriberModel) Language() *Language {
	return &t.language
}

// Speech-to-text request
// lang: "en"
// TODO: Transcriber does not have transcription stream, set it up when it does
func (o *TranscriberModel) STTRequest(r io.Reader) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		Language: "en",
		Reader: r,
	}
	resp, err := o.client.CreateTranscription(o.ctx, req)
	if err != nil {
		log.Printf("Transcription error: %v\n", err)
		return "", err
	}
	return resp.Text, nil
}


func (t *TranscriberModel) WriteRTP(pkt *rtp.Packet) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.oggSerializer == nil {
		oggSerializer, err := oggwriter.NewWith(t.oggWriter, t.rtpCodec.ClockRate, t.rtpCodec.Channels)
		if err != nil {
			println("failed to create ogg serializer", err)
			return err
		}
		t.oggSerializer = oggSerializer
	}

	//t.sb.Push(pkt)
	//for _, p := range t.sb.PopPackets() {
	if err := t.oggSerializer.WriteRTP(pkt); err != nil {
		return err
	}
	//}

	return nil
}

func (t *TranscriberModel) start() error {
	defer func() {
		close(t.closeCh)
	}()

	for {
		stream, err := t.newStream()
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Canceled {
				return nil
			}

			println("failed to create a new speech stream", err)
			t.results <- TranscriptionResult{
				Error: err,
			}
			return err
		}

		endStreamCh := make(chan struct{})
		nextCh := make(chan struct{})

		// Forward oggreader to the speech stream
		go func() {
			defer close(nextCh)
			buf := make([]byte, 1024)
			for {
				select {
				case <-endStreamCh:
					return
				default:
					n, err := t.oggReader.Read(buf)
					if err != nil {
						if err != io.EOF {
							println("failed to read from ogg reader", err)
						}
						return
					}

					if n <= 0 {
						continue // No data
					}

					if err := stream.Send(&sttpb.StreamingRecognizeRequest{
						StreamingRequest: &sttpb.StreamingRecognizeRequest_AudioContent{
							AudioContent: buf[:n],
						},
					}); err != nil {
						if err != io.EOF {
							println("failed to forward audio data to speech stream", err)
							t.results <- TranscriptionResult{
								Error: err,
							}
						}
						return
					}
				}
			}

		}()

		// Read transcription results
		for {
			resp, err := stream.Recv()
			if err != nil {
				if status, ok := status.FromError(err); ok {
					if status.Code() == codes.OutOfRange {
						break // Create a new speech stream (maximum speech length exceeded)
					} else if status.Code() == codes.Canceled {
						return nil // Context canceled (Stop)
					}
				}

				println("failed to receive response from speech stream", err)
				t.results <- TranscriptionResult{
					Error: err,
				}

				return err
			}

			if resp.Error != nil {
				break
			}

			// Read the whole transcription and put inside one string
			// We don't need to process each part individually (atm?)
			var sb strings.Builder
			final := false
			for _, result := range resp.Results {
				alt := result.Alternatives[0]
				text := alt.Transcript
				sb.WriteString(text)

				if result.IsFinal {
					sb.Reset()
					sb.WriteString(text)
					final = true
					break
				}
			}

			t.results <- TranscriptionResult{
				Text:    sb.String(),
				IsFinal: final,
			}
		}

		close(endStreamCh)

		// When nothing is written on the transcriber (The track is muted), this will block because the oggReader
		// is waiting for data. It avoids to create useless speech streams. (Also we end up here because Google automatically close the
		// previous stream when there's no "activity")
		//
		// Otherwise (When we have data) it is used to wait for the end of the current stream,
		// so we can create the next one and reset the oggSerializer
		<-nextCh

		// Create a new oggSerializer each time we open a new SpeechStream
		// This is required because the stream requires ogg headers to be sent again
		t.lock.Lock()
		t.oggSerializer = nil
		t.lock.Unlock()
	}
}

func (t *TranscriberModel) Close() {
	t.cancel()
	t.oggReader.Close()
	t.oggWriter.Close()
	<-t.closeCh
	close(t.results)
}

func (t *TranscriberModel) newStream() (sttpb.Speech_StreamingRecognizeClient, error) {
	stream, err := t.speechClient.StreamingRecognize(t.ctx)
	if err != nil {
		return nil, err
	}

	var helloClassItems []*sttpb.CustomClass_ClassItem
	var botClassItems []*sttpb.CustomClass_ClassItem
	for idx, v := range GreetingWords  {
		helloClassItems[idx].Value = v
	}
	for idx, v := range NameWords  {
		botClassItems[idx].Value = v
	}

	config := &sttpb.RecognitionConfig{
		Model: "command_and_search",
		Adaptation: &sttpb.SpeechAdaptation{
			PhraseSets: []*sttpb.PhraseSet{
				{
					Phrases: []*sttpb.PhraseSet_Phrase{
						{Value: "${hello} ${bot}"},
						{Value: "${bot}"},
						{Value: "Hey ${bot}"},
						{Value: "Kitt"},
						{Value: "Kit-t"},
						{Value: "Kit"},
					},
					Boost: 16,
				},
			},
			CustomClasses: []*sttpb.CustomClass{
				{
					CustomClassId: "hello",
					Items: helloClassItems,
				},
				{
					CustomClassId: "bot",
					Items: botClassItems,
				},
			},
		},
		UseEnhanced:       true,
		Encoding:          sttpb.RecognitionConfig_OGG_OPUS,
		SampleRateHertz:   int32(t.rtpCodec.ClockRate),
		AudioChannelCount: int32(t.rtpCodec.Channels),
		LanguageCode:      t.language.TranscriberCode,
	}

	if err := stream.Send(&sttpb.StreamingRecognizeRequest{
		StreamingRequest: &sttpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &sttpb.StreamingRecognitionConfig{
				InterimResults: true,
				Config:         config,
			},
		},
	}); err != nil {
		return nil, err
	}

	return stream, nil
}


func (t *TranscriberModel) Results() <-chan TranscriptionResult {
	return t.results
}
