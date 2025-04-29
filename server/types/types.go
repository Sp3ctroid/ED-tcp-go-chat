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

type JSON_payload struct {
	Username string `json:"author"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	Status   string `json:"status"`
}

type Room struct {
	Name  string
	Users map[string]*Client
}

func NewRoom(name string) *Room {
	return &Room{
		Name:  name,
		Users: make(map[string]*Client),
	}
}

type Message struct {
	Text   string
	Author *Client
	Time   string
	Dest   *Room
}

func NewMessage() *Message {
	return &Message{
		Text:   "",
		Author: nil,
		Dest:   nil,
	}
}

func (msg *Message) FillMessage(c *Client, text string) {
	msg.Text = text
	msg.Author = c
	msg.Dest = c.Room
	msg.Time = time.Now().Format(time.TimeOnly)
}

type Client struct {
	Username   string
	Reader     bufio.Reader
	Writer     bufio.Writer
	Connection net.Conn
	Room       *Room
}

func NewClient(connection net.Conn, server *Server) *Client {
	cl := &Client{
		Username:   "anon",
		Reader:     *bufio.NewReader(connection),
		Writer:     *bufio.NewWriter(connection),
		Connection: connection,
		Room:       nil,
	}
	server.InfoLog.Printf("New CLIENT: %s: %v", cl.Username, connection.RemoteAddr())
	go cl.Read(server)

	return cl

}

func (client *Client) Read(server *Server) {

	for {
		str, err := client.Reader.ReadString('\n')
		if err != nil {
			server.ErrorLog.Println("Reading from CLIENT: ", client.Username)
			return
		}

		msg := NewMessage()
		msg.FillMessage(client, str)

		server.Incoming <- msg
	}

}

func (client *Client) Write(Status string, Username string, Text string, Time string, server *Server) {

	JSON_payload := JSON_payload{Username: Username, Text: Text, Time: Time, Status: Status}
	str_json, err := json.Marshal(JSON_payload)

	_, err = client.Writer.WriteString(string(str_json))
	if err != nil {
		return
	}
	server.InfoLog.Printf("GOT PAYLOAD: %s, %s, %s, %s", JSON_payload.Username, JSON_payload.Text, JSON_payload.Time, JSON_payload.Status)
	server.InfoLog.Printf("Sent MESSAGE: (%s) to CLIENT: (%s); in ROOM: (%s)", strings.Trim(Text, "\n"), client.Username, client.Room.Name)
	client.Writer.Flush()

}

var cmd = map[string]int{
	"join":   1,
	"leave":  2,
	"list":   3,
	"help":   4,
	"create": 5,
}

const (
	CMDJoin   = 1
	CMDLeave  = 2
	CMDList   = 3
	CMDHelp   = 4
	CMDCreate = 5
)

type Server struct {
	Incoming   chan *Message
	Rooms      map[string]*Room
	Users      map[string]*Client
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	WarningLog *log.Logger
}

func NewServer(fileLog bool) *Server {

	file := os.Stdout
	err := error(nil)
	if fileLog {
		file, err = os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			panic(err)
		}
	}

	log.SetOutput(file)
	InfoLog := log.New(file, "[INFO] ", log.Ldate|log.Ltime)
	WarningLog := log.New(file, "[WARNING] ", log.Ldate|log.Ltime)
	ErrorLog := log.New(file, "[ERROR] ", log.Ldate|log.Ltime)

	InfoLog.Println("Server is starting...")

	s := &Server{
		Incoming:   make(chan *Message),
		Rooms:      make(map[string]*Room),
		Users:      make(map[string]*Client),
		InfoLog:    InfoLog,
		ErrorLog:   ErrorLog,
		WarningLog: WarningLog,
	}

	s.ListenChans()
	return s
}

func (s *Server) ListenChans() {
	go func() {
		for {
			select {
			case message_in := <-s.Incoming:
				s.InfoLog.Printf("Received MESSAGE: (%s) from CLIENT: (%s); in ROOM: (%s)", strings.Trim(message_in.Text, "\n"), message_in.Author.Username, message_in.Dest.Name)
				s.ParseMsg(message_in)
			}
		}
	}()
}

func (s *Server) ParseMsg(msg *Message) {
	switch {
	case msg.Text[0] == '/':
		msg.Text = msg.Text[1:]
		s.ParseCommand(msg)
	default:
		s.Broadcast(msg)
	}
}

func (s *Server) FormatText(msg Message) string {

	return fmt.Sprintf("%v %v: %v", msg.Time, msg.Author.Username, msg.Text)

}

func (s *Server) ParseCommand(msg *Message) {
	switch {
	case GetCommand(msg.Text) == CMDJoin:
		s.InfoLog.Println("Received JOIN command")
		s.JoinRoom(msg.Author, s.GetSecArg(msg))
	case GetCommand(msg.Text) == CMDLeave:
		s.InfoLog.Println("Received LEAVE command")
		s.LeaveRoom(msg.Author)
	case GetCommand(msg.Text) == CMDList:
		s.InfoLog.Println("Received LIST command")
		s.ListRooms(msg.Author)
	case GetCommand(msg.Text) == CMDHelp:
		s.InfoLog.Println("Received HELP command")
		//s.Help(msg.Author)
	case GetCommand(msg.Text) == CMDCreate:
		s.InfoLog.Println("Received CREATE command")
		s.CreateRoom(msg.Author, s.GetSecArg(msg))
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
	s.UtilMsgToClient(client, fmt.Sprintf("You left room: %s! You are now in room: General\n", client.Room.Name), time.Now().Format(time.TimeOnly), "LEFT", "System Notification")
	s.UtilBroadcast(client, " left room!\n", time.Now().Format(time.TimeOnly), "SENT", client.Username)
	s.InfoLog.Printf("CLIENT: (%s) left ROOM: (%s)", client.Username, client.Room.Name)

	delete(s.Rooms[client.Room.Name].Users, client.Username)
	client.Room = s.Rooms["General"]
	s.Rooms["General"].Users[client.Username] = client

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

func (s *Server) JoinRoom(client *Client, name string) {
	if s.Rooms[name] == nil {
		client.Write("NEX", "Server Notification", "Room Doesn't Exist\n", time.Now().Format(time.TimeOnly), s)
		s.WarningLog.Printf("CLIENT: (%s) tried to join ROOM: (%s) that does not exist", client.Username, name)
		return
	}

	delete(s.Rooms[client.Room.Name].Users, client.Username)
	client.Room = s.Rooms[name]
	s.Rooms[name].Users[client.Username] = client

	s.InfoLog.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room.Name)
	s.UtilMsgToClient(client, fmt.Sprintf("You joined room: %s!\n", name), time.Now().Format(time.TimeOnly), "JOINED", "System Notification")
	s.UtilBroadcast(client, "joined this room!\n", time.Now().Format(time.TimeOnly), "JOINED", client.Username)
}

func (s *Server) CreateRoom(client *Client, name string) {
	if s.Rooms[name] != nil {
		client.Write("ALREX", "Server Notification", "Room Already Exists\n", time.Now().Format(time.TimeOnly), s)
		s.WarningLog.Printf("CLIENT: (%s) tried to create ROOM: (%s) that already exists", client.Username, client.Room.Name)
		return
	}
	room := NewRoom(name)
	s.Rooms[name] = room
	delete(s.Rooms[client.Room.Name].Users, client.Username)
	client.Room = room
	room.Users[client.Username] = client
	s.InfoLog.Printf("CLIENT: (%s) created and joined ROOM: (%s)", client.Username, client.Room.Name)
	s.UtilMsgToClient(client, fmt.Sprintf("You Created and Joined room %s\n", client.Room.Name), time.Now().Format(time.TimeOnly), "CREATED\n", "System Notification")
	s.UtilBroadcastServer(s, "BRCREATED", name)
}

func (s *Server) ListRooms(client *Client) {
	s.InfoLog.Printf("CLIENT: (%s) requested list of rooms", client.Username)
	//s.UtilMsgToClient(client, "Available rooms:\n", )
	var room_list string
	idx := 1
	for _, room := range s.Rooms {
		if idx == len(s.Rooms) {
			room_list += strings.Trim(room.Name, "\n")
		} else {
			room_list += strings.Trim(room.Name, "\n") + " "
		}

		idx++

	}

	client.Write("LIST", "Server Notification", strings.Trim(room_list, "\n"), time.Now().Format(time.TimeOnly), s)
}

func (s *Server) UtilBroadcastServer(server *Server, Status string, Text string) {
	for _, user := range server.Users {
		user.Write(Status, "System Notification", Text, time.Now().Format(time.TimeOnly), s)
	}
}

func (s *Server) Broadcast(msg *Message) {
	for _, user := range msg.Dest.Users {
		user.Write("SENT", msg.Author.Username, msg.Text, msg.Time, s)
	}
}

func (s *Server) UtilBroadcast(client *Client, Text string, Time string, Status string, Sender string) {
	for _, user := range client.Room.Users {
		user.Write(Status, Sender, Text, Time, s)
	}
}
func (s *Server) UtilMsgToClient(client *Client, Text string, Time string, Status string, Sender string) {
	client.Write(Status, Sender, Text, Time, s)
}

func (s *Server) RecursiveUserNameCheck(client *Client) {
	if s.Users[client.Username] != nil {
		client.Username = client.Username + "*"
		s.RecursiveUserNameCheck(client)
	}
}

func (server *Server) Join(client *Client) {
	server.RecursiveUserNameCheck(client)
	server.Users[client.Username] = client

	server.Rooms["General"].Users[client.Username] = client
	client.Room = server.Rooms["General"]
	server.InfoLog.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room.Name)
}
