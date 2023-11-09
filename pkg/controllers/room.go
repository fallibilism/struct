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

	room_name := c.Query("room_name", "")
	if room_name == "" {
		return RoomIdError
	}

	rm := models.NewRoomModel(config.App)
	rm.Lock()
	defer rm.Unlock()
	if _, err := rm.GetRoomByName(room_name); err == nil {
		return RoomExistsError
	}
	if err := rm.CreateRoom(room_name, user_id);err != nil {
		return fiber.ErrBadGateway
	}
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
	log.Panicf("[%s] Not implemented", "HandleEndRoom")
	return nil
}
