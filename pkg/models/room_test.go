package models

import (
	"testing"

	"v/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestRoom(t *testing.T) {
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
	t.Run("test room", func(t *testing.T) {
		room := NewRoomModel(conf)
		err := room.CreateRoom(Room{
			RoomName: "bystander",
			IsActive: true,
		})

		assert.Equal(t, err, nil)
	})


	t.Run("test active rooms", func(t *testing.T) {
		room := NewRoomModel(conf)
		var rooms []*Room;
		rooms, err = room.GetActiveRooms()

		assert.Equal(t, err, nil)
		_, err = room.GetRoom(rooms[0].RoomId.String())
		assert.Equal(t, err, nil)
		_, err = room.GetRoomBySid(rooms[0].Sid)
		assert.Equal(t, err, nil)
	})
}
