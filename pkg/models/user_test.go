package models

import (
	"testing"
	"v/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	testDB, err := config.NewDbConnection(&config.PostgresConfig{
		Host: "localhost",
		Port: 5432,
	})

	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}

	conf := &config.AppConfig{
		DB: testDB,
	}
	t.Run("test user", func(t *testing.T) {
		user := NewUserModel(conf)
		err := user.Create("test", "user")

		assert.Equal(t, err, nil)
	})

	t.Run("test user validation", func(t *testing.T) {
		conf := config.TestConfig
		user := NewUserModel(conf)
		ctx := &fiber.Ctx{}
		ctx.Locals("admin", false)
		ctx.Locals("roomId", "")
		err := user.Validation(ctx)

		assert.Equal(t, err, fiber.NewError(fiber.StatusUnauthorized, AdminOnlyError))

		ctx.Locals("admin", true)
		err = user.Validation(ctx)

		assert.Equal(t, err, fiber.NewError(fiber.StatusNotFound, "Room ID not found"))

		ctx.Locals("roomId", "1")
		err = user.Validation(ctx)

		assert.Equal(t, err, nil)

	})
}
