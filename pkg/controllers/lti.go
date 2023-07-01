package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"v/pkg/config"
	"v/pkg/models"
	"v/pkg/utils"

	protocol "github.com/fallibilism/protocol/go_protocol"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func HandleLTIV1CheckRoom(c *fiber.Ctx) error {
	roomId := c.Locals("room_id")

	if roomId == nil {
		return fiber.NewError(fiber.StatusBadRequest, "room id is empty")
	}

	m := models.NewRoomAuthModel(config.App)
	status, err := m.IsRoomActive(&protocol.IsRoomActiveRequest{
		RoomId: roomId.(string),
	})

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

func HandleLTIV1JoinRoom(c *fiber.Ctx) error {
	return nil
}

func HandleLTIEndRoom(c *fiber.Ctx) error {
	roomId := c.Locals("room_id")
	isAdmin := c.Locals("is_admin").(bool)

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
func bearerSplit(tok string) (string, error) {
	if len(tok) > 6 && tok[:6] == "Bearer" {
		return tok[7:], nil
	}
	return "", errors.New("bearer token not found")
}
func HandleV1HeaderToken(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")

	if authToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, InvalidTokenError)
	}

	authToken, err := bearerSplit(authToken)

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, InvalidTokenError+": "+err.Error())
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

func HandleLTIAuth(c *fiber.Ctx) error {
	params := c.AllParams()

	if params == nil || len(params) < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "params is empty")
	}

	// url
	host := "http"
	if c.Secure() {
		host += "s"
	}

	if !isLocalhost(c.Hostname()) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid host. Only secured host is allowed")
	}

	url := fmt.Sprintf("%s://%s%s", host, c.Hostname(), c.Path())
	m := models.NewLTIV1Model(config.App)

	httpReq, err := adaptor.ConvertRequest(c, false)

	if err != nil {
		return err
	}

	ltis, err := m.LTIVerifyAuth(params, url, httpReq)
	if err != nil {
		return err
	}

	roomId := fmt.Sprintf("%s_%s_%s", ltis.Get("tool_consumer_instance_guid"), ltis.Get("context_id"), ltis.Get("resource_link_id"))
	userId := ltis.Get("user_id")
	if userId == "" {
		userId = utils.GenHash(ltis.Get("lis_person_contact_email_primary"))
	}

	if userId == "" {
		return errors.New("either value of user_id or lis_person_contact_email_primary  required")
	}

	name := ltis.Get("lis_person_name_full")
	if name == "" {
		name = "User_" + userId
	}

	claims := &protocol.LtiAuthClaims{
		UserId:    userId,
		Name:      name,
		IsAdmin:   false,
		RoomId:    utils.GenHash(roomId),
		RoomTitle: ltis.Get("context_label"),
	}

	if strings.Contains(ltis.Get("roles"), "Instructor") {
		claims.IsAdmin = true
	}

	tok, err := utils.ClaimsToJWT(claims)

	if err != nil {
		return err
	}

	vals := fiber.Map{
		"Title":   claims.RoomTitle,
		"Token":   tok,
		"IsAdmin": claims.IsAdmin,
	}

	if claims.LtiCustomParameters.LtiCustomDesign != nil {
		design, err := json.Marshal(claims.LtiCustomParameters.LtiCustomDesign)
		if err == nil {
			vals["CustomDesign"] = string(design)
		}

	}

	return c.Render("assets/lti/v1", vals)
}

func isLocalhost(host string) bool {
	return strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "0.0.0.0")
}
