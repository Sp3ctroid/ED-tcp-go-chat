package storage

import (
	"serverMod/types/rc"
	"strings"
)

func (store *ClientStoreMap) GET_User_By_Name(name string) *rc.Client {
	return store.Items[name]
}

type ClientStoreMap struct {
	Items map[string]*rc.Client
}

type RoomStoreMap struct {
	Items map[string]*Room
}

func NewRoomMap() *RoomStoreMap {
	return &RoomStoreMap{Items: make(map[string]*Room)}
}

func NewClientMap() *ClientStoreMap {
	return &ClientStoreMap{Items: make(map[string]*rc.Client)}
}

func (store *RoomStoreMap) DELETE_From_Room(roomName string, userName string) {
	delete(store.Items[roomName].Users.GET_All_Users_Server(), userName)
}

func (store *RoomStoreMap) ADD_To_Room(client *rc.Client, roomName string) {
	store.Items[roomName].Users.ADD_User_To_Server(client)
}

func (store *RoomStoreMap) ASSIGN_To_Room(client *rc.Client, roomName string) {
	client.Room = roomName
}

func (store *RoomStoreMap) CREATE_New_Room(room *Room) {
	store.Items[room.Name] = room
}

func (store *RoomStoreMap) CHECK_If_Exists(roomName string) bool {
	return store.Items[roomName] != nil
}

func (store *RoomStoreMap) GET_All_Rooms() string {
	var room_list string
	idx := 1
	for _, room := range store.Items {
		if idx == len(store.Items) {
			room_list += strings.Trim(room.Name, "\n")
		} else {
			room_list += strings.Trim(room.Name, "\n") + " "
		}

		idx++

	}

	return room_list
}

func (store *RoomStoreMap) GET_Room(roomName string) *Room {
	return store.Items[roomName]
}

func (store *ClientStoreMap) GET_All_Users_Server() map[string]*rc.Client {
	return store.Items
}

func (store *RoomStoreMap) GET_All_Users_Room(roomName string) map[string]*rc.Client {
	return store.Items[roomName].Users.GET_All_Users_Server()
}

func (store *ClientStoreMap) CHECK_If_Exists(name string) bool {
	return store.Items[name] != nil
}

func (store *ClientStoreMap) UPDATE_Username(name string, changeTo string) {

	client := store.Items[name]
	client.Username = changeTo
	store.Items[changeTo] = client
	delete(store.Items, name)
}

func (store *ClientStoreMap) ADD_User_To_Server(client *rc.Client) {
	store.Items[client.Username] = client
}
