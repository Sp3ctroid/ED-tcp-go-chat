package storage

import (
	"serverMod/types/rc"
)

type Room struct {
	Name  string
	Users UserStore
}
type UserStore interface {
	GET_All_Users_Server() map[string]*rc.Client
	CHECK_If_Exists(name string) bool
	UPDATE_Username(name string, changeTo string)
	ADD_User_To_Server(client *rc.Client)
	GET_User_By_Name(name string) *rc.Client
}

type RoomStore interface {
	DELETE_From_Room(roomName string, userName string)
	ADD_To_Room(client *rc.Client, roomName string)
	ASSIGN_To_Room(client *rc.Client, roomName string)
	CREATE_New_Room(room *Room)
	CHECK_If_Exists(roomName string) bool
	GET_All_Rooms() string
	GET_All_Users_Room(roomName string) map[string]*rc.Client
	GET_Room(roomName string) *Room
}

func NewRoom(name string) *Room {
	return &Room{
		Name:  name,
		Users: NewClientMap(),
	}
}
