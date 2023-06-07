package models

import (
	"database/sql"
	"errors"
	"v/pkg/config"

	"github.com/gofiber/fiber/v2"
)

type UserModel struct {
	db *sql.DB
	rs *RoomService
}

type User struct {
	RoomId   string `json:"room_id"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

func NewUserModel() *UserModel {
	return &UserModel{
		db: config.App.DB,
		rs: NewRoomService(),
	}
}

// TODO: 😴 replace with gorm
func (u *UserModel) Create(user *User) error {
	_, err := u.db.Exec("INSERT INTO users (id, room_id, name, role, is_active) VALUES (?, ?, ?, ?, ?)", user.Id, user.RoomId, user.Name, user.Role, user.IsActive)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserModel) Validation(c *fiber.Ctx) error {
	admin := c.Locals("admin")
	roomId := c.Params("roomId")

	if isAdmin, ok := admin.(bool); ok {
		if !isAdmin {
			return fiber.NewError(fiber.StatusUnauthorized, "You are not authorized to access this room")
		}
	}

	if roomId == "" {
		return fiber.NewError(fiber.StatusNotFound, "Room not found")
	}

	return nil
}

func (u *UserModel) SwitchPresenter() error {
	return errors.New("[SwitchPresenter] implement me")
}
