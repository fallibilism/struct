package models

import (
	"context"
	"errors"
	"log"
	"v/pkg/config"
	protocol "v/protocol/go_protocol"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
	"gorm.io/gorm"

)

/**
* This file is part of the v project.
* it is responsible for the room services in **livekit**.
* referenced from [github.com/livekit-examples](github.com/livekit-examples)
*
 */
type RoomService struct {
	app           *config.AppConfig
	db            *gorm.DB
	livekitClient *lksdk.RoomServiceClient
	ctx           context.Context
}

var (
	errorRoomNotFound = errors.New("room not found")
	errorRoomExists   = errors.New("room already exists")
)

const (
	errorRoomNotActive = "room is not active"
)

// a livekit room service
func NewRoomService(conf *config.AppConfig) *RoomService {
	livekitClient := lksdk.NewRoomServiceClient(config.Livekit.Host, config.Livekit.ApiKey, config.Livekit.Secret)
	return &RoomService{
		db:            conf.DB,
		livekitClient: livekitClient,
		ctx:           context.Background(),
	}
}

// load all livekit rooms given a room id
func (r *RoomService) LoadRoom(roomId string) (*livekit.Room, error) {
	req := livekit.ListRoomsRequest{
		Names: []string{
			roomId,
		},
	}

	rooms, err := r.livekitClient.ListRooms(r.ctx, &req)
	if err != nil {
		log.Printf("Error loading room: %v", err)
		return nil, err
	}

	if len(rooms.Rooms) == 0 {
		return nil, errorRoomNotFound
	}

	return rooms.Rooms[0], nil
}


func (rm *RoomService) JoinRoom(user_id string, url string, admin bool, room *livekit.Room) (*Room, error) {
	var r Room

	if u := rm.db.Model(&Room{}).Preload("User").Where("room_id = ?", room.Name).First(&r).Error; u != nil {
		return nil, RoomDoesNotExistError
	}
	println("participants: [", r.Participants , "]")

	if len(r.Participants) > 0 {
		println("bot already connected[room name: ", room.Name, ", count: ", room.NumParticipants)
		return nil, errors.New("bot already connected")
	}
	u := NewUserModel(rm.app)

	if err := u.ChangeState(user_id, protocol.ConnectionState_CONNECTING); err != nil { 
		return nil, err
	}
	
	tm := NewTokenModel(rm.app)
	token,err := tm.AddToken(room.Sid, config.BotIdentity, admin)
	if err != nil {
		return nil, err
	}

	ConnectGPTParticipant(config.App, token, url)
	println("Bot connected successfully, token: ", token)
	return nil, nil


}
// load livekit room participants given a room id
func (r *RoomService) LoadParticipants(roomId string) ([]*livekit.ParticipantInfo, error) {
	req := livekit.ListParticipantsRequest{
		Room: roomId,
	}

	participants, err := r.livekitClient.ListParticipants(context.Background(), &req)
	if err != nil {
		log.Printf("Error loading participants: %v", err)
		return nil, err
	}

	return participants.Participants, nil
}

func (r *RoomService) CreateRoom(roomId string) (*livekit.Room, error) {

	if _, err := r.LoadRoom(roomId); err == nil {
		if err == errorRoomNotFound {
			return nil, errorRoomExists
		}
		return nil, err
	}

	req := livekit.CreateRoomRequest{
		Name: roomId,
	}

	room, err := r.livekitClient.CreateRoom(r.ctx, &req)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		return nil, err
	}

	return room, nil
}

// delete livekit room given a room id
func (r *RoomService) DeleteRoom(roomId string) (error) {
	req := livekit.DeleteRoomRequest{
		Room: roomId,
	}

	_, err := r.livekitClient.DeleteRoom(r.ctx, &req)
	if err != nil {
		log.Printf("Error deleting room: %v", err)
		return err
	}

	return nil
}
