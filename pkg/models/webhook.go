package models

import (
	"context"
	"log"
	"v/pkg/config"

	"github.com/livekit/protocol/livekit"
	"github.com/livekit/protocol/webhook"
	"github.com/redis/go-redis/v9"
)

const (
	errorUnknownEvent = "unknown event"
)

type webhookEvent struct {
	rc          *redis.Client
	ctx         context.Context
	event       *livekit.WebhookEvent
	RoomService *RoomService
	// notifier       *WebhookNotifierModel
	// userModel      *UserModel
	// recoderModel   *RecorderModel
	// recordingModel *RecordingModel
	// roomModel      *RoomModel
}

func NewWebhook(conf *config.AppConfig, event *livekit.WebhookEvent) {
	wh := &webhookEvent{
		rc:          conf.Redis,
		ctx:         context.Background(),
		event:       event,
		RoomService: NewRoomService(conf),
	}

	switch event.GetEvent() {
	case webhook.EventRoomStarted:
		wh.roomStarted()
	case webhook.EventRoomFinished:
		wh.roomFinished()
	case webhook.EventParticipantJoined:
		wh.participantJoined()
	case webhook.EventParticipantLeft:
		wh.participantLeft()
	case webhook.EventTrackPublished:
		wh.trackPublished()
	case webhook.EventTrackUnpublished:
		wh.trackUnpublished()
	default:
		log.Printf("%s: %s", errorUnknownEvent, event.GetEvent())
	}

}

func (wh *webhookEvent) roomStarted() {
	log.Printf("room started: %s", wh.event.Room.Name)

	panic("implement me")
}

func (wh *webhookEvent) roomFinished() {
	log.Printf("room finished: %s", wh.event.Room.Name)

	panic("implement me")
}

func (wh *webhookEvent) participantJoined() {
	log.Printf("participant joined: %s", wh.event.Participant.Identity)

	panic("implement me")
}

func (wh *webhookEvent) participantLeft() {
	log.Printf("participant left: %s", wh.event.Participant.Identity)

	panic("implement me")
}

func (wh *webhookEvent) trackPublished() {
	log.Printf("track published: %s", wh.event.Track.Sid)
}

func (wh *webhookEvent) trackUnpublished() {
	log.Printf("track unpublished: %s", wh.event.Track.Sid)
}
