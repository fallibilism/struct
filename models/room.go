package models

import (
	"context"
	"database/sql"
	"sync"
	"time"
	"v/pkg/config"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID
	RoomName     string    `json:"room_name"`
	RoomId       string    `json:"room_id"`
	Created      string    `json:"created"`
	CreatedAt    time.Time `json:"created_at"`
	IsRunning    bool      `json:"is_running"`
	IsRecording  bool      `json:"is_recording"`
	WebhookUrl   string    `json:"webhook_url"`
	RecorderId   string    `json:"recorder_id"`
	IsActiveRTMP bool      `json:"is_active_rtmp"`
	Ended        bool      `json:"ended"`
}

type RoomModel struct {
	app *config.AppConfig
	db  *sql.DB
	ctx context.Context
	sync.Mutex
}

func NewRoomModel() *RoomModel {
	return &RoomModel{
		app: config.App,
		db:  config.App.DB,
		ctx: context.Background(),
	}
}
