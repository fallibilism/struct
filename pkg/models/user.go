package models

import (
	"errors"
	"sync"
	"time"
	"v/pkg/config"
	protocol "v/protocol/go_protocol"


	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

var (
	InactiveState protocol.ConnectionState = *protocol.ConnectionState_INACTIVE.Enum()
	ActiveState protocol.ConnectionState = *protocol.ConnectionState_ACTIVE.Enum()
	ConnectiongState protocol.ConnectionState = *protocol.ConnectionState_CONNECTING.Enum()
	ErrUserDoesNotExist error = errors.New("the user does not exist")
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
	State  protocol.ConnectionState `gorm:"type:uuid;default:0"`
	Name     string `gorm:"not null"`
	Role     string 
	IsActive bool   
}

func (u *User) BeforeCreate(*gorm.DB) error {
	u.ID = shortuuid.New()
	u.IsActive = true
	return nil
}



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
		IsActive: true,
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
func (u *UserModel) ChangeState(user_id string, state protocol.ConnectionState) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	if err := u.db.Model(&User{}).Where("id = ?", user_id).Update("state = ?", state).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserModel) Get(name string) *User {
	var user *User
	if err := u.db.Where("name = ?", name).First(user).Error; err != nil {
		return nil 
	}
	return user
}

func (u *UserModel) Validation(user_id, room_id string, admin bool) bool {
	var user *User
	u.lock.Lock()
	defer u.lock.Unlock()
	if err := u.db.Where("id = ? & room_id = ? & is_admin = ?", user_id, room_id, admin).First(user).Error; err != nil {
		return false
	}
	if u == nil {
		return false
	}
	return true
}

func (u *UserModel) SwitchPresenter() error {
	return errors.New("[SwitchPresenter] implement me")
}
