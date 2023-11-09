package models

import (
	"context"
	"errors"
	"sync"
	"time"
	"v/pkg/config"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/livekit/protocol/livekit"
	"gorm.io/gorm"
)

/*
 * Model for room in DB
 * This is the model for the room table in the database
 * It is used to store the room information
 * It is also used to store the room information in the database
 */
type Room struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
//	Sid          string    `gorm:""`//set later
	RoomName     string    `gorm:"not null,unique"`
	RoomId       string `gorm:"unique"`
	UserId      string    `gorm:""` // created by who?
	Participants []User `gorm:"foreignKey:RoomId;references:RoomId"`
	IsActive     bool      
	IsRecording  bool      `gorm:""`
	RecorderId   string    `gorm:""`
	WebhookUrl   string    `gorm:""`
	IsActiveRTMP bool      `gorm:""`
	Ended		time.Time 
}

func (r *Room) BeforeCreate(*gorm.DB) error {
	r.RoomId = shortuuid.New()
	r.IsActive = false
	r.IsRecording = false
	r.IsActiveRTMP = false
	return nil
}

var (
	RoomDoesNotExistError = errors.New("room does not exist")
)

type RoomModel struct {
	app *config.AppConfig
	db  *gorm.DB
	ctx context.Context
	lock sync.Mutex
}

func NewRoomModel(conf *config.AppConfig) *RoomModel {
	return &RoomModel{
		app: conf,
		db:  conf.DB,
		ctx: context.Background(),
	}
}

// Get info about the room from db
func (rm *RoomModel) GetRoom(roomId string) (*Room, error) {
	println("roomid: ", roomId)
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

// Get info about the room from db by room name
func (rm *RoomModel) GetRoomByName(room_name string) (*Room, error) {
	println("roomid: ", room_name)
	uid, err := uuid.Parse(room_name)

	if err != nil {
		return nil, err
	}

	var room Room
	if err := rm.db.Where("room_name = ?", uid).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (rm *RoomModel) JoinRoom(user_id string, room *livekit.Room) (*Room, error) {
	rm.lock.Lock()
	var r Room

	if u := rm.db.Model(&Room{}).Preload("User").Where("room_id = ?", room.Metadata).First(&r).Error; u != nil {
		return nil, RoomDoesNotExistError
	}
	println("participants: [", r.Participants , "]")

	if len(r.Participants) > 0 {
		rm.lock.Unlock()
		println("bot already connected[room name: ", room.Name, ", count: ", room.NumParticipants)
		return nil, errors.New("bot already connected")
	}
	u := NewUserModel(rm.app)

	if err := u.ChangeState(user_id, ConnectingState); err != nil { 
		return nil, err
	}
	rm.lock.Unlock()
	
	
	tm := NewTokenModel(rm.app)
	token,err := tm.AddToken(room.Sid, config.BotIdentity)
	if err != nil {
		return nil, err
	}

	println("Bot connected successfully, token: ", token)
	return nil, nil


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
func (rm *RoomModel) CreateRoom(room_name string, user_id string) error {
	if err := rm.db.Create(&Room{
		RoomName: room_name,
		UserId: user_id,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (rm *RoomModel) UpdateRoom(room *Room) error {
	var r []*Room
	if err := rm.db.Where("room_id = ?", room.RoomId).Find(&r).Error; err != nil {
		return err
	}

	if len(r) > 0 {
		if err := rm.db.Where("room_id = ?", room.RoomId).Updates(room).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}
