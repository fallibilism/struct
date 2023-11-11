package models

import (
	"errors"
	"sync"
	"time"
	"v/pkg/config"
	"v/protocol"

	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

// user's current state
type State string

const (
	ActiveState State = "ACTIVE"
	ConnectingState State = "CONNECTING"
	InactiveState State = "INACTIVE"
	x protocol.State = 3
)

type UserModel struct {
	db *gorm.DB
	rs *RoomService
	lock sync.Mutex
}

type User struct {
	ID        string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RoomId string 
	State  
	Name     string `gorm:"not null"`
	Role     string 
	IsActive bool   
}

func (u *User) BeforeCreate(*gorm.DB) error {
	u.ID = shortuuid.New()
	u.IsActive = true
	return nil
}


const (
	AdminOnlyError = "You are not authorized to access this room. Only admin can access this room"
)

func NewUserModel(conf *config.AppConfig) *UserModel {
	return &UserModel{
		db: conf.DB,
		rs: NewRoomService(conf),
	}
}

func (u *UserModel) Create(name, role string) error {
	if err := u.db.Create(&User{
		Name:	name,
		Role:	role,
	}).Error; err != nil {
		return err
	}

	return nil
}

// add user to a room
func (u *UserModel) AddRoom(user_id, room_id string) error {
	if err := u.db.Model(&User{}).Where("id = ?", user_id).Update("room_id = ?", room_id).Error; err != nil {
		return err
	}
	return nil
}

// alternate between state ["ACTIVE", "CONNECTING", "INACTIVE"]
func (u *UserModel) ChangeState(user_id string, state State) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	if err := u.db.Model(&User{}).Where("id = ?", user_id).Update("state = ?", state).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Validation(c *fiber.Ctx) error {
	admin := c.Locals("admin")
	roomId := c.Locals("roomId")

	if isAdmin, ok := admin.(bool); ok && !isAdmin {
		if !isAdmin {
			return fiber.NewError(fiber.StatusUnauthorized, AdminOnlyError)
		}
	}

	if roomId == "" {
		return fiber.NewError(fiber.StatusNotFound, "Room ID not found")
	}

	return nil
}

func (u *UserModel) SwitchPresenter() error {
	return errors.New("[SwitchPresenter] implement me")
}
