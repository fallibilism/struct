package models

import (
	"context"
	"log"
	"time"
	"v/pkg/config"

	"github.com/livekit/protocol/livekit"
	livekitClient "github.com/livekit/server-sdk-go"

	protocol "github.com/fallibilism/protocol/go_protocol"
)

type RoomAuthModel struct {
	rm *RoomModel
	rs *RoomService
}

func NewRoomAuthModel(conf *config.AppConfig) *RoomAuthModel {
	return &RoomAuthModel{
		rm: NewRoomModel(conf),
		rs: NewRoomService(conf),
	}
}

// check if room is active
func (m *RoomAuthModel) IsRoomActive(req *protocol.IsRoomActiveRequest) (bool, error) {
	r, err := m.rm.GetRoom(req.RoomId)
	if err != nil {
		return false, err
	}

	if r.ID < 1 {
		return false, errorRoomNotFound
	}
	room, err := m.LoadRoom(req.RoomId)
	if err != nil {
		return false, err
	}

	return room.State == livekit.ParticipantInfo_JOINED, nil
}

// load all livekit rooms given a room id
func (m *RoomAuthModel) LoadRoom(roomId string) (*livekit.Room, error) {
	req := livekit.ListRoomsRequest{
		Names: []string{
			roomId,
		},
	}

	rooms, err := livekitClient.ListRooms(context.Background(), &req)
	if err != nil {
		log.Printf("Error loading room: %v", err)
		return nil, err
	}

	if len(rooms.Rooms) == 0 {
		return nil, errorRoomNotFound
	}

	return rooms.Rooms[0], nil
}

// end a livekit room given a room id
func (m *RoomAuthModel) EndRoom(req *protocol.EndRoomRequest) error {
	m.rm.Lock()
	defer m.rm.Unlock()
	room, err := m.rm.GetRoom(req.RoomId)

	if err != nil {
		return err
	}

	if room == nil {
		return errorRoomNotFound
	}

	_, err = m.rs.DeleteRoom(req.RoomId)

	if err != nil {
		return err
	}

	m.rm.UpdateRoom(&Room{
		RoomId:   room.RoomId,
		IsActive: false,
		Ended:    time.Now().Format("2008-06-03 15:04:05"),
	})

}
