package models

import (
	"errors"
	"v/pkg/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserModel struct {
	db *gorm.DB
	rs *RoomService
}

type User struct {
	RoomId   string `json:"room_id"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
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

func (u *UserModel) Create(user *User) error {
	if err := u.db.Create(&User{
		Id:       user.Id,
		RoomId:   user.RoomId,
		Name:     user.Name,
		Role:     user.Role,
		IsActive: user.IsActive,
	}).Error; err != nil {
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
