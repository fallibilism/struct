package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func HandleAuthHeaderCheck(c *fiber.Ctx) error {
	api_key := c.Get("API_KEY", "")
	hash_signature := c.Get("HASH", "")
	body := c.Body()

	// check if api key exists

	// unhash signature

	// compare signatures

	log.Panicf("[%s] Not implemented", "HandleAuthHeaderCheck")
	return c.Next()
	//return nil
}

func HandleGenerateJoinToken(c *fiber.Ctx) error {
	log.Panicf("[%s] Not implemented", "HandleGenerateJoinToken")
	return nil
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
