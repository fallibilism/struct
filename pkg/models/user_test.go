package models

import (
	"testing"
	"v/pkg/config"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func migrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&Room{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Token{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

func TestUser(t *testing.T) {
	testDB, err := config.NewDbConnection(&config.PostgresConfig{
		Host: "localhost",
		Port: 5432,
	})

	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}

	err = migrations(testDB)
	if err != nil {
		t.Fatalf("Cannot connect to database: %v", err)
	}
	

	t.Run("test user", func(t *testing.T) {
		conf := &config.AppConfig{
			DB: testDB,
		}
		user := NewUserModel(conf)
		err := user.Create("test", "user")
		assert.Equal(t, err, nil)

		u1, err := user.Get("test")
		assert.Equal(t, err, nil)

		err = user.Create("test2", "admin")
		assert.Equal(t, err, nil)

		u2, err := user.Get("test2")
		assert.Equal(t, err, nil)

		t.Run("add user to room", func(t *testing.T) {

			err = user.AddRoom(u1.ID, "room1")
			assert.Equal(t, err, nil)

			err = user.AddRoom(u2.ID, "room1")
			assert.Equal(t, err, nil)
		})

		t.Run("validate user", func(t *testing.T) {
			v := user.Validation(u1.ID, u1.RoomId, u1.Role)
			assert.Equal(t, v, nil, "%v", *u1)



			v = user.Validation(u2.ID, u2.RoomId, u2.Role)
			assert.Equal(t, v, nil, "%v", *u2)
		})

	})

}
