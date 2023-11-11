package controllers

import (
	"log"
	"v/pkg/config"
	"v/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// -- with auth room handlers --

func HandleRoomCreate(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	user_id,ok := claims["user_id"].(string)

	if !ok {
		return fiber.ErrUnauthorized
	}

	room_id := c.Query("room_id", "")
	admin := c.Query("admin", "false") == "true"
	if room_id == "" {
		return RoomIdError
	}

	rm := models.NewRoomModel(config.App)
	if _, err := rm.GetRoomByName(room_id); err != nil {
		if err == models.RoomDoesNotExistError {
			if err := rm.CreateRoom(room_id, user_id);err != nil {
				return fiber.ErrBadGateway
			}
		}
	}

	rs := models.NewRoomService(config.App)
	room, err := rs.LoadRoom(room_id)
	if err != nil {
		return RoomIdError
	}
	rs.JoinRoom(user_id, config.Conf.Livekit.Host, admin, room)
	return c.JSON("OK")

}

func HandleRoomActivity(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandlerRoomActivity")
	return nil
}

func HandleActiveRoomInfo(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandleActiveRoomInfo")
	return nil
}
func HandleActiveRoomsInfo(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandleActiveRoomsInfo")
	return nil
}

func HandleEndRoom(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	user_id,ok := claims["user_id"].(string)

	if !ok {
		return fiber.ErrUnauthorized
	}

	room_id := c.Query("room_id", "")
	if room_id == "" {
		return RoomIdError
	}

	admin := c.Query("admin", "false") == "true"
	if !admin {
		return fiber.ErrUnauthorized
	}

	um := models.NewUserModel(config.App)
	if um.Validation(user_id, room_id, admin) {
		return fiber.ErrBadRequest
	}

	rm := models.NewRoomModel(config.App)
	if _, err := rm.DeleteRoom(room_id); err != nil {
		return fiber.ErrInternalServerError
	}

	rs := models.NewRoomService(config.App)
	err := rs.DeleteRoom(room_id)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON("OK")
}
