package models

import (
	"context"
	"sync"
	"time"
	"v/pkg/config"

	"github.com/lithammer/shortuuid/v4"
	"github.com/livekit/protocol/auth"
	"gorm.io/gorm"
)

type Token struct {
	ID string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RoomId string `json:"room_id"`
	UserID   string `json:"user_id"` // id of user who created it
	Token	string `json:"token"`
}
func (r *Token) BeforeCreate(*gorm.DB) error {
	r.ID = shortuuid.New()
	return nil
}

type TokenModel struct {
	db  *gorm.DB
	ctx context.Context
	config *config.Config
	sync.Mutex
}

func NewTokenModel(app *config.AppConfig) *TokenModel {
	return &TokenModel{
		db:  app.DB,
		ctx: context.Background(),
	}
}

func (t *TokenModel) AddToken(room_id, user_id string) (string, error) {
	if err := t.db.Where("room_id = ?", room_id).First(&Room{}).Error; err != nil {
		return "", RoomDoesNotExistError
	}

	lk_api_key, lk_secret := config.Conf.Livekit.ApiKey,config.Conf.Livekit.Secret
	at := auth.NewAccessToken(lk_api_key, lk_secret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room_id,
	}
	at.AddGrant(grant).
		SetIdentity(user_id).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return "", err
	}

	tc := Token{Token: token, UserID: user_id}
	if err := t.db.Where("room_id = ?", room_id).Attrs(tc).FirstOrCreate(&tc).Error; err != nil {
		return "", err
	}

	return token, nil
}
