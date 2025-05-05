package types

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
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

type JSON_payload struct {
	Username string `json:"author"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	Status   string `json:"status"`
}

type Room struct {
	Name  string
	Users UserStore
}

type UserStore interface {
	GET_All_Users_Server() map[string]*Client
	CHECK_If_Exists(name string) bool
	UPDATE_Username(name string, changeTo string)
	ADD_User_To_Server(client *Client)
	GET_User_By_Name(name string) *Client
}

type RoomStore interface {
	DELETE_From_Room(roomName string, userName string)
	ADD_To_Room(client *Client, roomName string)
	ASSIGN_To_Room(client *Client, roomName string)
	CREATE_New_Room(room *Room)
	CHECK_If_Exists(roomName string) bool
	GET_All_Rooms() string
	GET_All_Users_Room(roomName string) map[string]*Client
}

func NewRoom(name string) *Room {
	return &Room{
		Name:  name,
		Users: NewClientMap(),
	}
}

type Message struct {
	Text   string
	Author string
	Time   string
	Dest   string
}

func NewMessage() *Message {
	return &Message{
		Text:   "",
		Author: "",
		Dest:   "",
	}
}

func (msg *Message) FillMessage(c *Client, text string) {
	msg.Text = text
	msg.Author = c.Username
	msg.Dest = c.Room
	msg.Time = time.Now().Format(time.TimeOnly)
}

type Client struct {
	Username   string
	Reader     bufio.Reader
	Writer     bufio.Writer
	Connection net.Conn
	Room       string
}

func NewClient(connection net.Conn, server *Server) *Client {
	cl := &Client{
		Username:   "anon",
		Reader:     *bufio.NewReader(connection),
		Writer:     *bufio.NewWriter(connection),
		Connection: connection,
		Room:       "",
	}
	INFOLOG.Printf("New CLIENT: %s: %v", cl.Username, connection.RemoteAddr())
	go cl.Read(server)

	return cl

}

func (client *Client) Read(server *Server) {

	for {
		str, err := client.Reader.ReadString('\n')
		if err != nil {
			ERRORLOG.Println("Reading from CLIENT: ", client.Username)
			return
		}

		msg := NewMessage()
		msg.FillMessage(client, str)

		server.Incoming <- msg
	}

}

func (client *Client) Write(Status string, Username string, Text string, Time string) {

	JSON_payload := JSON_payload{Username: Username, Text: Text, Time: Time, Status: Status}
	str_json, err := json.Marshal(JSON_payload)

	if err != nil {
		return
	}

	_, err = client.Writer.WriteString(string(str_json))
	if err != nil {
		return
	}
	INFOLOG.Printf("GOT PAYLOAD: %s, %s, %s, %s", JSON_payload.Username, JSON_payload.Text, JSON_payload.Time, JSON_payload.Status)
	INFOLOG.Printf("Sent MESSAGE: (%s) to CLIENT: (%s); in ROOM: (%s)", strings.Trim(Text, "\n"), client.Username, client.Room)
	client.Writer.Flush()

}

type Server struct {
	Incoming chan *Message
	Rooms    RoomStore
	Users    UserStore
}

var STREAM = os.Stdout
var INFOLOG = log.New(STREAM, "[INFO] ", log.Ldate|log.Ltime)
var WARNINGLOG = log.New(STREAM, "[WARNING] ", log.Ldate|log.Ltime)
var ERRORLOG = log.New(STREAM, "[ERROR] ", log.Ldate|log.Ltime)

func NewServer(fileLog bool) *Server {

	if fileLog {
		file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			panic(err)
		}

		log.SetOutput(file)
	}

	INFOLOG.Println("Server is starting...")

	s := &Server{
		Incoming: make(chan *Message),
		Rooms:    NewRoomMap(),
		Users:    NewClientMap(),
	}

	s.ListenChans()
	return s
}

func (s *Server) ListenChans() {
	go func() {
		for {
			select {
			case message_in := <-s.Incoming:
				INFOLOG.Printf("Received MESSAGE: (%s) from CLIENT: (%s); in ROOM: (%s)", strings.Trim(message_in.Text, "\n"), message_in.Author, message_in.Dest)
				s.ParseMsg(message_in)
			}
		}
	}()
}

func (store *ClientStoreMap) GET_User_By_Name(name string) *Client {
	return store.Items[name]
}

func (s *Server) ParseMsg(msg *Message) {
	switch {
	case msg.Text[0] == '/':
		msg.Text = msg.Text[1:]
		s.ParseCommand(msg, s.Users.GET_User_By_Name(msg.Author))
	default:
		s.Broadcast(msg)
	}
}

func (s *Server) FormatText(msg Message) string {

	return fmt.Sprintf("%v %v: %v", msg.Time, msg.Author, msg.Text)

}

func (s *Server) ParseCommand(msg *Message, clinet *Client) {
	switch {
	case GetCommand(msg.Text) == CMDJoin:
		INFOLOG.Println("Received JOIN command")
		s.JoinRoom(clinet, s.GetSecArg(msg))
	case GetCommand(msg.Text) == CMDLeave:
		INFOLOG.Println("Received LEAVE command")
		s.LeaveRoom(clinet)
	case GetCommand(msg.Text) == CMDList:
		INFOLOG.Println("Received LIST command")
		s.ListRooms(clinet)
	case GetCommand(msg.Text) == CMDHelp:
		INFOLOG.Println("Received HELP command")
		//s.Help(msg.Author)
	case GetCommand(msg.Text) == CMDCreate:
		INFOLOG.Println("Received CREATE command")
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

func (s *Server) GetSecArg(msg *Message) string {
	return strings.Trim(strings.Split(msg.Text, " ")[1], "\n")
}

func (s *Server) LeaveRoom(client *Client) {
	s.UtilMsgToClient(client, fmt.Sprintf("You left room: %s! You are now in room: General\n", client.Room), time.Now().Format(time.TimeOnly), "LEFT", "System Notification")
	s.UtilBroadcast(client, " left room!\n", time.Now().Format(time.TimeOnly), "SENT", client.Username)
	INFOLOG.Printf("CLIENT: (%s) left ROOM: (%s)", client.Username, client.Room)

	s.Rooms.DELETE_From_Room(client.Room, client.Username) //DELETE FROM ROOM
	s.Rooms.ASSIGN_To_Room(client, "General")              //ASSIGN CLIENT ROOM
	s.Rooms.ADD_To_Room(client, "General")                 //ADD TO ROOM

}

type ClientStoreMap struct {
	Items map[string]*Client
}

type RoomStoreMap struct {
	Items map[string]*Room
}

func NewRoomMap() *RoomStoreMap {
	return &RoomStoreMap{Items: make(map[string]*Room)}
}

func NewClientMap() *ClientStoreMap {
	return &ClientStoreMap{Items: make(map[string]*Client)}
}

func (store *RoomStoreMap) DELETE_From_Room(roomName string, userName string) {
	delete(store.Items[roomName].Users.GET_All_Users_Server(), userName)
}

func (store *RoomStoreMap) ADD_To_Room(client *Client, roomName string) {
	store.Items[roomName].Users.ADD_User_To_Server(client)
}

func (store *RoomStoreMap) ASSIGN_To_Room(client *Client, roomName string) {
	client.Room = roomName
}

func (store *RoomStoreMap) CREATE_New_Room(room *Room) {
	store.Items[room.Name] = room
}

func (store *RoomStoreMap) CHECK_If_Exists(roomName string) bool {
	return store.Items[roomName] != nil
}

func (s *Server) JoinRoom(client *Client, name string) {
	if !s.Rooms.CHECK_If_Exists(name) { //CHECK IF EXISTS
		client.Write("NEX", "Server Notification", "Room Doesn't Exist\n", time.Now().Format(time.TimeOnly))
		WARNINGLOG.Printf("CLIENT: (%s) tried to join ROOM: (%s) that does not exist", client.Username, name)
		return
	}

	s.UtilBroadcast(client, "left this room!\n", time.Now().Format(time.TimeOnly), "LEFT", client.Username)

	s.Rooms.DELETE_From_Room(client.Room, client.Username) //DELETE FROM ROOM
	s.Rooms.ASSIGN_To_Room(client, name)                   //ASSIGN TO ROOM
	s.Rooms.ADD_To_Room(client, name)                      //ADD TO ROOM

	INFOLOG.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room)
	s.UtilMsgToClient(client, fmt.Sprintf("You joined room: %s!\n", name), time.Now().Format(time.TimeOnly), "JOINED", "System Notification")
	s.UtilBroadcast(client, "joined this room!\n", time.Now().Format(time.TimeOnly), "JOINED", client.Username)
}

func (s *Server) CreateRoom(client *Client, name string) {
	if s.Rooms.CHECK_If_Exists(name) { //CHECK IF EXISTS
		client.Write("ALREX", "Server Notification", "Room Already Exists\n", time.Now().Format(time.TimeOnly))
		WARNINGLOG.Printf("CLIENT: (%s) tried to create ROOM: (%s) that already exists", client.Username, client.Room)
		return
	}

	s.UtilBroadcast(client, "left this room!\n", time.Now().Format(time.TimeOnly), "LEFT", client.Username)

	room := NewRoom(name)
	s.Rooms.CREATE_New_Room(room)                          //CREATE NEW ROOM
	s.Rooms.DELETE_From_Room(client.Room, client.Username) //DELETE FROM ROOM
	s.Rooms.ASSIGN_To_Room(client, room.Name)              //ASSIGN TO ROOM
	s.Rooms.ADD_To_Room(client, room.Name)                 //ADD TO ROOM
	INFOLOG.Printf("CLIENT: (%s) created and joined ROOM: (%s)", client.Username, client.Room)
	s.UtilMsgToClient(client, fmt.Sprintf("You Created and Joined room %s\n", client.Room), time.Now().Format(time.TimeOnly), "CREATED\n", "System Notification")
	s.UtilBroadcastServer(s, "BRCREATED", name)
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

func (s *Server) ChangeUserName(oldName, newUserName string) {
	cl := s.Users.GET_User_By_Name(oldName)
	if s.Users.CHECK_If_Exists(newUserName) {

		cl.Write("ALRTAK", "System Notification", "Username is already taken", time.Now().Format(time.TimeOnly))
		return
	}

	s.UtilBroadcast(s.Users.GET_User_By_Name(oldName), fmt.Sprintf("changed his name to: %s", newUserName), time.Now().Format(time.TimeOnly), "USERNAMECHANGED", oldName)

	s.Rooms.DELETE_From_Room(s.Users.GET_User_By_Name(oldName).Room, oldName)
	s.Users.UPDATE_Username(oldName, newUserName)
	s.Rooms.ADD_To_Room(s.Users.GET_User_By_Name(newUserName), s.Users.GET_User_By_Name(newUserName).Room)

	s.UtilMsgToClient(cl, fmt.Sprintf("You changed your username to %s", newUserName), time.Now().Format(time.TimeOnly), "CHANGED", "System Notification")

}

func (s *Server) ListRooms(client *Client) {
	INFOLOG.Printf("CLIENT: (%s) requested list of rooms", client.Username)
	//s.UtilMsgToClient(client, "Available rooms:\n", )
	room_list := s.Rooms.GET_All_Rooms()

	client.Write("LIST", "Server Notification", strings.Trim(room_list, "\n"), time.Now().Format(time.TimeOnly))
}

func (store *ClientStoreMap) GET_All_Users_Server() map[string]*Client {
	return store.Items
}

func (s *Server) UtilBroadcastServer(server *Server, Status string, Text string) {
	for _, user := range s.Users.GET_All_Users_Server() { // GET ALL USERS SERVER-WIDE
		user.Write(Status, "System Notification", Text, time.Now().Format(time.TimeOnly))
	}
}

func (store *RoomStoreMap) GET_All_Users_Room(roomName string) map[string]*Client {
	return store.Items[roomName].Users.GET_All_Users_Server()
}

func (s *Server) Broadcast(msg *Message) {
	roomName := msg.Dest
	for _, user := range s.Rooms.GET_All_Users_Room(roomName) { //GET ALL USERS IN A ROOM
		user.Write("SENT", msg.Author, msg.Text, msg.Time)
	}
}

func (s *Server) UtilBroadcast(client *Client, Text string, Time string, Status string, Sender string) {
	roomName := client.Room
	for _, user := range s.Rooms.GET_All_Users_Room(roomName) { //GET ALL USERS IN A ROOM
		user.Write(Status, Sender, Text, Time)
	}
}
func (s *Server) UtilMsgToClient(client *Client, Text string, Time string, Status string, Sender string) {
	client.Write(Status, Sender, Text, Time)
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

func (s *Server) RecursiveUserNameCheck(client *Client) {
	if s.Users.CHECK_If_Exists(client.Username) { //CHECK IF USER EXISTS
		client.Username = client.Username + "*" //UPDATE NAME
		s.RecursiveUserNameCheck(client)
	}
}

func (store *ClientStoreMap) ADD_User_To_Server(client *Client) {
	store.Items[client.Username] = client
}

func (server *Server) Join(client *Client) {
	server.RecursiveUserNameCheck(client)
	server.Users.ADD_User_To_Server(client) //ADD USER TO SERVER STORAGE

	server.Rooms.ADD_To_Room(client, "General")    //ADD TO ROOM
	server.Rooms.ASSIGN_To_Room(client, "General") //ASSIGN ROOM
	INFOLOG.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room)
}

// func (s *Server) Help(client *Client) {
// 	client.Write("Available commands:\n")
// 	client.Write("\n")
// 	client.Write("join <name>\n")
// 	client.Write("\n")
// 	client.Write("leave\n")
// 	client.Write("\n")
// 	client.Write("list\n")
// 	client.Write("\n")
// 	client.Write("create <name>\n")
// 	client.Write("\n")
// }
