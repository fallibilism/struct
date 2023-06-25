package models

import (
	"context"
	"sync"
	"v/pkg/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID       uuid.UUID
	RoomName string `json:"room_name"`
	RoomId   string `json:"room_id"`
	Created  string `json:"created"`
	UserID   string `json:"user_id"`
}

type TokenModel struct {
	app *config.AppConfig
	db  *gorm.DB
	ctx context.Context
	sync.Mutex
}

func NewTokenModel(conf *config.AppConfig) *TokenModel {
	return &TokenModel{
		db:  conf.DB,
		ctx: context.Background(),
	}
}
