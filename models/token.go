package models

import (
	"context"
	"database/sql"
	"sync"
	"v/pkg/config"

	"github.com/google/uuid"
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
	db  *sql.DB
	ctx context.Context
	sync.Mutex
}

func NewTokenModel() *TokenModel {
	return &TokenModel{
		app: config.App,
		db:  config.App.DB,
		ctx: context.Background(),
	}
}
