package models

import (
	"context"
	"sync"
	"time"
	"v/pkg/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	db  *gorm.DB
	ctx context.Context
	sync.Mutex
}

func NewRoomModel(conf *config.AppConfig) *RoomModel {
	return &RoomModel{
		db:  conf.DB,
		ctx: context.Background(),
	}
}
