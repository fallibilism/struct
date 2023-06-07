package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"v/pkg/config"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"

	"github.com/redis/go-redis/v9"
)

type RoomService struct {
	app           *config.AppConfig
	db            *sql.DB
	redis         *redis.Client
	livekitClient *lksdk.RoomServiceClient
	ctx           context.Context
}

var (
	errorRoomNotFound = errors.New("room not found")
	errorRoomExists   = errors.New("room already exists")
)

func NewRoomService() *RoomService {
	livekitClient := lksdk.NewRoomServiceClient(config.Livekit.Host, config.Livekit.ApiKey, config.Livekit.Secret)
	return &RoomService{
		app:           config.App,
		db:            config.App.DB,
		redis:         config.App.Redis,
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

func (r *RoomService) DeleteRoom(roomId string) (string, err error) {
	req := livekit.DeleteRoomRequest{
		Room: roomId,
	}

	res, err := r.livekitClient.DeleteRoom(r.ctx, &req)
	x := res
	if err != nil {
		log.Printf("Error deleting room: %v", err)
		return "", err
	}

	return 
}
