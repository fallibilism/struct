package models

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"v/pkg/config"

	"github.com/jordic/lti"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	AuthError = "error authenticating user"
)

type LtiV1 struct {
	authModel      *RoomAuthModel
	authTokenModel *TokenModel
}

type LtiClaims struct {
	UserId              string               `json:"user_id"`
	Name                string               `json:"name"`
	IsAdmin             bool                 `json:"is_admin"`
	RoomId              string               `json:"room_id"`
	RoomTitle           string               `json:"room_title"`
	LtiCustomParameters *LtiCustomParameters `json:"lti_custom_parameters,omitempty"`
}

type LtiCustomParameters struct {
	RoomDuration               uint64           `json:"room_duration,omitempty"`
	AllowPolls                 *bool            `json:"allow_polls,omitempty"`
	AllowSharedNotePad         *bool            `json:"allow_shared_note_pad,omitempty"`
	AllowBreakoutRoom          *bool            `json:"allow_breakout_room,omitempty"`
	AllowRecording             *bool            `json:"allow_recording,omitempty"`
	AllowRTMP                  *bool            `json:"allow_rtmp,omitempty"`
	AllowViewOtherWebcams      *bool            `json:"allow_view_other_webcams,omitempty"`
	AllowViewOtherParticipants *bool            `json:"allow_view_other_users_list,omitempty"`
	MuteOnStart                *bool            `json:"mute_on_start,omitempty"`
	LtiCustomDesign            *LtiCustomDesign `json:"lti_custom_design,omitempty"`
}

type LtiCustomDesign struct {
	PrimaryColor    string `json:"primary_color,omitempty"`
	SecondaryColor  string `json:"secondary_color,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	CustomLogo      string `json:"custom_logo,omitempty"`
}

func NewLTIV1Model(conf *config.AppConfig) *LtiV1 {
	return &LtiV1{
		authModel:      NewRoomAuthModel(conf),
		authTokenModel: NewTokenModel(conf),
	}
}

func (l *LtiV1) LTIVerifyHeaderToken(token string) (*LtiClaims, error) {
	tok, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	out := jwt.Claims{}
	claims := &LtiClaims{}
	if err = tok.Claims([]byte(config.Conf.JWTSecret), &out, claims); err != nil {
		return nil, err
	}
	if err = out.Validate(jwt.Expected{Issuer: config.Conf.JWTIssuer, Time: time.Now()}); err != nil {
		return nil, err
	}

	return claims, nil

}

// this verifies the LTI 1.1 authentication request
func (l *LtiV1) LTIVerifyAuth(body map[string]string, url string) (*url.Values, error) {
	p := lti.NewProvider(config.Conf.JWTSecret, url)
	p.Method = "POST"
	p.ConsumerKey = config.Conf.ConsumerKey

	// mapParams, providedSignature := split(body)

	//---
	providedSignature := body["oauth_signature"]
	delete(body, "oauth_signature")
	//---

	for k, v := range body {
		p.Add(k, v)
	}

	if p.Get("oauth_consumer_key") != p.ConsumerKey {
		return nil, fmt.Errorf("%s : invalid consumer_key", AuthError)
	}

	sign, err := p.Sign()
	if err != nil {
		return nil, fmt.Errorf(" %s : %s", AuthError, err.Error())
	}

	params := p.Params()
	if sign != providedSignature {
		err := fmt.Errorf("Expected: " + sign + "but provided: " + providedSignature)
		log.Println(err)
		return nil, errors.New(AuthError + ": verification failed")
	}

	return &params, nil
}

// uncomment when post request is fixed
// // splits params to maps
// func split(par string) (m map[string]string, providedSig string) {
// 	req := strings.Split(par, "&")
// 	m = make(map[string]string)

// 	for _, v := range req {
// 		s := strings.Split(v, "=")
// 		b, _ := url.QueryUnescape(s[1])
// 		m[s[0]] = s[1]
// 		if s[0] == "oauth_signature" {
// 			providedSig = b
// 		}
// 	}

// 	return m, providedSig
// }
