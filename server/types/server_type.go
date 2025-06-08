package types

import (
	"fmt"
	"log"
	"os"
	"serverMod/types/logger"
	"serverMod/types/rc"
	"serverMod/types/storage"
	"strings"
	"time"
)

var cmd = map[string]int{
	"join":   1,
	"leave":  2,
	"list":   3,
	"help":   4,
	"create": 5,
	"name":   6,
}

const (
	CMDJoin       = 1
	CMDLeave      = 2
	CMDList       = 3
	CMDHelp       = 4
	CMDCreate     = 5
	CMDChangeName = 6
)

type Server struct {
	Incoming chan *rc.Message
	Rooms    storage.RoomStore
	Users    storage.UserStore
}

func NewServer(fileLog bool) *Server {

	if fileLog {
		file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			panic(err)
		}

		log.SetOutput(file)
	}

	logger.INFOLOG.Println("Server is starting...")

	s := &Server{
		Incoming: make(chan *rc.Message),
		Rooms:    storage.NewRoomMap(),
		Users:    storage.NewClientMap(),
	}

	s.ListenChans()
	return s
}

func (s *Server) ListenChans() {
	go func() {
		for {
			select {
			case message_in := <-s.Incoming:
				logger.INFOLOG.Printf("Received MESSAGE: (%s) from CLIENT: (%s); in ROOM: (%s)", strings.Trim(message_in.Text, "\n"), message_in.Author, message_in.Dest)
				s.ParseMsg(message_in)
			}
		}
	}()
}

func (s *Server) ParseMsg(msg *rc.Message) {
	switch {
	case msg.Text[0] == '/':
		msg.Text = msg.Text[1:]
		s.ParseCommand(msg, s.Users.GET_User_By_Name(msg.Author))
	default:
		s.Broadcast(msg)
	}
}

func (s *Server) FormatText(msg rc.Message) string {

	return fmt.Sprintf("%v %v: %v", msg.Time, msg.Author, msg.Text)

}

func (s *Server) ParseCommand(msg *rc.Message, clinet *rc.Client) {
	switch {
	case GetCommand(msg.Text) == CMDJoin:
		logger.INFOLOG.Println("Received JOIN command")
		s.JoinRoom(clinet, s.GetSecArg(msg))
	case GetCommand(msg.Text) == CMDLeave:
		logger.INFOLOG.Println("Received LEAVE command")
		s.LeaveRoom(clinet)
	case GetCommand(msg.Text) == CMDList:
		logger.INFOLOG.Println("Received LIST command")
		s.ListRooms(clinet)
	case GetCommand(msg.Text) == CMDHelp:
		logger.INFOLOG.Println("Received HELP command")
		//s.Help(msg.Author)
	case GetCommand(msg.Text) == CMDCreate:
		logger.INFOLOG.Println("Received CREATE command")
		s.CreateRoom(clinet, s.GetSecArg(msg))
	case GetCommand(msg.Text) == CMDChangeName:
		s.ChangeUserName(msg.Author, s.GetSecArg(msg))
	}
}

func GetCommand(text string) int {
	command := strings.Split(text, " ")[0]
	command = strings.Trim(command, "\n")
	return cmd[command]
}

func (s *Server) GetSecArg(msg *rc.Message) string {
	return strings.Trim(strings.Split(msg.Text, " ")[1], "\n")
}

func (s *Server) LeaveRoom(client *rc.Client) {
	s.UtilMsgToClient(client, fmt.Sprintf("You left room: %s! You are now in room: General\n", client.Room), time.Now().Format(time.TimeOnly), "LEFT", "System Notification")
	s.UtilBroadcast(client, " left room!\n", time.Now().Format(time.TimeOnly), "SENT", client.Username)
	logger.INFOLOG.Printf("CLIENT: (%s) left ROOM: (%s)", client.Username, client.Room)

	s.Rooms.DELETE_From_Room(client.Room, client.Username) //DELETE FROM ROOM
	s.Rooms.ASSIGN_To_Room(client, "General")              //ASSIGN CLIENT ROOM
	s.Rooms.ADD_To_Room(client, "General")                 //ADD TO ROOM

}

func (s *Server) JoinRoom(client *rc.Client, name string) {
	if !s.Rooms.CHECK_If_Exists(name) { //CHECK IF EXISTS
		client.Write("NEX", "Server Notification", "Room Doesn't Exist\n", time.Now().Format(time.TimeOnly))
		logger.WARNINGLOG.Printf("CLIENT: (%s) tried to join ROOM: (%s) that does not exist", client.Username, name)
		return
	}

	s.UtilBroadcast(client, "left this room!\n", time.Now().Format(time.TimeOnly), "LEFT", client.Username)

	s.Rooms.DELETE_From_Room(client.Room, client.Username) //DELETE FROM ROOM
	s.Rooms.ASSIGN_To_Room(client, name)                   //ASSIGN TO ROOM
	s.Rooms.ADD_To_Room(client, name)                      //ADD TO ROOM

	logger.INFOLOG.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room)
	s.UtilMsgToClient(client, fmt.Sprintf("You joined room: %s!\n", name), time.Now().Format(time.TimeOnly), "JOINED", "System Notification")
	s.UtilBroadcast(client, "joined this room!\n", time.Now().Format(time.TimeOnly), "USERJOINED", client.Username)
}

func (s *Server) CreateRoom(client *rc.Client, name string) {
	if s.Rooms.CHECK_If_Exists(name) { //CHECK IF EXISTS
		client.Write("ALREX", "Server Notification", "Room Already Exists\n", time.Now().Format(time.TimeOnly))
		logger.WARNINGLOG.Printf("CLIENT: (%s) tried to create ROOM: (%s) that already exists", client.Username, client.Room)
		return
	}

	s.UtilBroadcast(client, "left this room!\n", time.Now().Format(time.TimeOnly), "LEFT", client.Username)

	room := storage.NewRoom(name)
	s.Rooms.CREATE_New_Room(room)                          //CREATE NEW ROOM
	s.Rooms.DELETE_From_Room(client.Room, client.Username) //DELETE FROM ROOM
	s.Rooms.ASSIGN_To_Room(client, room.Name)              //ASSIGN TO ROOM
	s.Rooms.ADD_To_Room(client, room.Name)                 //ADD TO ROOM
	logger.INFOLOG.Printf("CLIENT: (%s) created and joined ROOM: (%s)", client.Username, client.Room)
	s.UtilMsgToClient(client, fmt.Sprintf("You Created and Joined room %s\n", client.Room), time.Now().Format(time.TimeOnly), "CREATED\n", "System Notification")
	s.UtilBroadcastServer(s, "BRCREATED", name)
}

func (s *Server) ChangeUserName(oldName, newUserName string) {
	cl := s.Users.GET_User_By_Name(oldName)
	if s.Users.CHECK_If_Exists(newUserName) {

		cl.Write("ALRTAK", "System Notification", "Username is already taken", time.Now().Format(time.TimeOnly))
		return
	}

	s.UtilBroadcast(s.Users.GET_User_By_Name(oldName), fmt.Sprintf("changed his name to: %s", newUserName), time.Now().Format(time.TimeOnly), "USERNAMECHANGED", oldName)

	s.Rooms.GET_Room(s.Users.GET_User_By_Name(oldName).Room).Users.UPDATE_Username(oldName, newUserName)
	s.Users.UPDATE_Username(oldName, newUserName)

	s.UtilMsgToClient(cl, fmt.Sprintf("You changed your username to %s", newUserName), time.Now().Format(time.TimeOnly), "CHANGED", "System Notification")

}

func (s *Server) ListRooms(client *rc.Client) {
	logger.INFOLOG.Printf("CLIENT: (%s) requested list of rooms", client.Username)
	//s.UtilMsgToClient(client, "Available rooms:\n", )
	room_list := s.Rooms.GET_All_Rooms()

	client.Write("LIST", "Server Notification", strings.Trim(room_list, "\n"), time.Now().Format(time.TimeOnly))
}

func (s *Server) UtilBroadcastServer(server *Server, Status string, Text string) {
	for _, user := range s.Users.GET_All_Users_Server() { // GET ALL USERS SERVER-WIDE
		user.Write(Status, "System Notification", Text, time.Now().Format(time.TimeOnly))
	}
}

func (s *Server) Broadcast(msg *rc.Message) {
	roomName := msg.Dest
	for _, user := range s.Rooms.GET_All_Users_Room(roomName) { //GET ALL USERS IN A ROOM
		user.Write("SENT", msg.Author, msg.Text, msg.Time)
	}
}

func (s *Server) UtilBroadcast(client *rc.Client, Text string, Time string, Status string, Sender string) {
	roomName := client.Room
	for _, user := range s.Rooms.GET_All_Users_Room(roomName) { //GET ALL USERS IN A ROOM
		user.Write(Status, Sender, Text, Time)
	}
}

func (s *Server) UtilMsgToClient(client *rc.Client, Text string, Time string, Status string, Sender string) {
	client.Write(Status, Sender, Text, Time)
}

func (s *Server) RecursiveUserNameCheck(client *rc.Client) {
	if s.Users.CHECK_If_Exists(client.Username) { //CHECK IF USER EXISTS
		client.Username = client.Username + "*" //UPDATE NAME
		s.RecursiveUserNameCheck(client)
	}
}

func (server *Server) Join(client *rc.Client) {
	server.RecursiveUserNameCheck(client)
	server.Users.ADD_User_To_Server(client) //ADD USER TO SERVER STORAGE

	server.Rooms.ADD_To_Room(client, "General")    //ADD TO ROOM
	server.Rooms.ASSIGN_To_Room(client, "General") //ASSIGN ROOM
	logger.INFOLOG.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room)
}
