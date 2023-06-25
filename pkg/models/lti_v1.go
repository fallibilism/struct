package models

import (
	"time"

	"v/pkg/config"

	"gopkg.in/square/go-jose.v2/jwt"
)

type LtiV1 struct {
	authModel      *RoomAuthModel
	authTokenModel *AuthTokenModel
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

func NewLTIV1Model() *LtiV1 {
	return &LtiV1{
		authModel:      NewRoomAuthModel(),
		authTokenModel: NewAuthTokenModel(),
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
