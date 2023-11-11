package controllers

import (
	"log"
	"v/pkg/config"
	"v/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	InvalidTokenError = "The Token in the Authorization header is invalid"

)

var (
	RoomIdError = fiber.NewError(fiber.ErrBadRequest.Code, "room id is invalid")
	RoomExistsError = fiber.NewError(fiber.ErrBadRequest.Code, "room exists")
)
func HandleAuthHeaderCheck(c *fiber.Ctx) error {
	api_key := c.Get("API_KEY", "")
	hash_signature := c.Get("HASH", "")
	body := c.Body()

	println(api_key)
	println(hash_signature)
	println(body)
	// check if api key exists

	// unhash signature

	// compare signatures

	log.Panicf("[%s] Not implemented", "HandleAuthHeaderCheck")
	return c.Next()
	//return nil
}

// to regenerate join token
func HandleGenerateJoinToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	user_id, ok := claims["user_id"].(string)
	if !ok {
		return fiber.ErrBadRequest
	}

	room_id := c.Query("room_id")
	admin := c.Query("admin") == "true"


	if room_id == "" {
		return RoomIdError
	}

	tm := models.NewTokenModel(config.App)
	tm.Lock()
	defer tm.Unlock()

	t, err := tm.AddToken(room_id, user_id, admin)
	if err != nil {
		return err	
	}

	return c.JSON(t)
}

func HandleVerifyHeaderToken(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")

	if authToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// validate token

	log.Panicf("[%s] Not implemented", "HandleVerifyHeaderToken")
	return c.Next()
}

func HandleVerifyToken(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandleVerifyToken")
	return nil
}

func HandleRenewToken(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandleRenewToken")
	return nil
}

func HandleRevokeToken(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandleRevokeToken")
	return nil
}
