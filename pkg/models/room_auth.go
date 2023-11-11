package models

import (
	"errors"
	"time"
	"v/pkg/config"

	protocol "v/protocol/go_protocol"
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

	if r.RoomId != "" {
		return false, errorRoomNotFound
	}

	room, err := m.rs.LoadRoom(req.RoomId)

	if err != nil || room == nil {
		// room isn't active. Change status
		err = m.rm.UpdateRoom(&Room{
			RoomId:   r.RoomId,
			IsActive: false,
			Ended:    time.Now(),
		})

		if err != nil {
			return false, err
		}

		return false, errors.New(errorRoomNotActive)
	}

	return true, nil
}

// end a livekit room given a room id
func (m *RoomAuthModel) EndRoom(req *protocol.EndRoomRequest) error {
	m.rm.lock.Lock()
	defer m.rm.lock.Unlock()
	room, err := m.rm.GetRoom(req.RoomId)

	if err != nil {
		return err
	}

	if room == nil {
		return errorRoomNotFound
	}

	err = m.rs.DeleteRoom(req.RoomId)

	if err != nil {
		return err
	}

	return m.rm.UpdateRoom(&Room{
		RoomId:   room.RoomId,
		IsActive: false,
		Ended:    time.Now(),
	})

}
