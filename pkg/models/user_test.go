package models

import (
	"testing"
	"v/pkg/config"

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

		err = user.Create("test2", "admin")
		assert.Equal(t, err, nil)
	})

	t.Run("test user validation", func(t *testing.T) {
		conf := config.TestConfig
		user := NewUserModel(conf)
		test := user.Get("test")
		test2 := user.Get("test2")
		assert.NotNil(t, test)
		assert.NotNil(t, test2)


		v := user.Validation(test.ID, test.RoomId, test.Role == "admin")
		assert.Equal(t, v, false)

		v = user.Validation(test2.ID, test2.RoomId, test2.Role == "admin")
		assert.Equal(t, v, true)

		v = user.Validation(test.ID, test.RoomId, test.Role == "user")
		assert.Equal(t, v, true)


	})
}
