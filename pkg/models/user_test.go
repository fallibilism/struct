package models

import (
	"testing"

	"v/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	testDB, err := config.NewMockDbConnection()

	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}

	err = testDB.AutoMigrate(&Room{})

	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}

	err = testDB.AutoMigrate(&User{})

	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}
	conf := &config.AppConfig{
		DB: testDB,
	}
	t.Run("test user", func(t *testing.T) {
		user := NewUserModel(conf)
		err := user.Create(&User{
			Id:       "1",
			RoomId:   "1",
			Name:     "test",
			Role:     "admin",
			IsActive: true,
		})

		assert.Equal(t, err, nil)
	})
}
