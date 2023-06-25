package controllers

import (
	"encoding/json"
	"v/pkg/config"
	"v/pkg/models"

	protocol "github.com/fallibilism/protocol/go_protocol"
	"github.com/gofiber/fiber/v2"
)

func HandleLTIV1CheckRoom(c *fiber.Ctx) error {
	roomId := c.Locals("roomId")

	m := models.NewRoomAuthModel(config.App)
	status, msg := m.IsRoomActive(&protocol.IsRoomActiveRequest{
		RoomId: roomId.(string),
	})

	return c.JSON(fiber.Map{
		"status": status,
		"msg":    msg,
	})
}

func HandleLTIV1JoinRoom(c *fiber.Ctx) error {
	return nil
}

func HandleLTIEndRoom(c *fiber.Ctx) error {
	roomId := c.Locals("roomId")
	isAdmin := c.Locals("isAdmin").(bool)

	if !isAdmin {
		return c.JSON(fiber.Map{
			"status": false,
			"msg":    "only admin can perform this",
		})
	}

	m := models.NewRoomAuthModel(config.App)
	err := m.EndRoom(&protocol.EndRoomRequest{
		RoomId: roomId.(string),
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"msg":    "room ended successfully",
	})
}
func HandleV1HeaderToken(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")

	if authToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, InvalidTokenError)
	}

	m := models.NewLTIV1Model(config.App)
	claims, err := m.LTIVerifyHeaderToken(authToken)

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, InvalidTokenError+" "+err.Error())
	}

	c.Locals("room_id", claims.RoomId)
	c.Locals("room_title", claims.RoomTitle)
	c.Locals("user_id", claims.UserId)
	c.Locals("name", claims.Name)
	c.Locals("is_admin", claims.IsAdmin)

	if claims.LtiCustomParameters != nil {
		c.Locals("room_duration", claims.LtiCustomParameters.RoomDuration)
		customParams, err := json.Marshal(claims.LtiCustomParameters)
		if err == nil && customParams != nil {
			c.Locals("custom_params", customParams)
		}
	}

	return c.Next()
}
