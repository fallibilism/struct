package models

import (
	"context"
	"sync"
	"v/pkg/config"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
 * Model for room in DB
 * This is the model for the room table in the database
 * It is used to store the room information
 * It is also used to store the room information in the database
 */
type Room struct {
	*gorm.Model
	Sid          string    `json:"sid"`
	RoomName     string    `json:"room_name"`
	RoomId       uuid.UUID `json:"room_id"`
	Created      string    `json:"created"` // created by who?
	IsActive     bool      `json:"is_active"`
	IsRecording  bool      `json:"is_recording"`
	WebhookUrl   string    `json:"webhook_url"`
	RecorderId   string    `json:"recorder_id"`
	IsActiveRTMP bool      `json:"is_active_rtmp"`
	Ended        string    `json:"ended"`
	// CreatedAt    time.Time `json:"created_at"` // given by gorm.Model
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

// Get info about the room from db
func (rm *RoomModel) GetRoom(roomId string) (*Room, error) {
	uid, err := uuid.Parse(roomId)

	if err != nil {
		return nil, err
	}

	var room Room
	if err := rm.db.Where("room_id = ?", uid).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// Get info about the room from db by sid
func (rm *RoomModel) GetRoomBySid(sid string) (*Room, error) {
	var room Room
	if err := rm.db.Where("sid = ?", sid).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// Get info about the all active rooms from db
func (rm *RoomModel) GetActiveRooms() ([]*Room, error) {
	var rooms []*Room
	if err := rm.db.Where("is_active = ?", true).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// Create or update room in db
func (rm *RoomModel) CreateOrUpdateRoom(room *Room) error {
	if err := rm.db.Save(room).Error; err != nil {
		return err
	}
	return nil
}

func (rm *RoomModel) UpdateRoom(room *Room) error {
	var r []*Room
	if err := rm.db.Where("room_id = ?", room.RoomId).Find(&r).Error; err != nil {
		return err
	}

	if len(r) > 1 {
		if err := rm.db.Where("room_id = ?", room.RoomId).Updates(room).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}
