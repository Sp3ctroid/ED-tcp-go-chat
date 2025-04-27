package main

import (
	"bufio"
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
}

const ()

type Server struct {
	Incoming   chan *Message
	Rooms      map[string]*Room
	Users      map[string]*Client
	infoLog    *log.Logger
	errorLog   *log.Logger
	warningLog *log.Logger
}

func NewServer() *Server {
	file, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	log.SetOutput(file)
	infoLog := log.New(file, "[INFO] ", log.Ldate|log.Ltime)
	warningLog := log.New(file, "[WARNING] ", log.Ldate|log.Ltime)
	errorLog := log.New(file, "[ERROR] ", log.Ldate|log.Ltime)

	infoLog.Println("Server is starting...")

	s := &Server{
		Incoming:   make(chan *Message),
		Rooms:      make(map[string]*Room),
		Users:      make(map[string]*Client),
		infoLog:    infoLog,
		errorLog:   errorLog,
		warningLog: warningLog,
	}

	s.ListenChans()
	return s
}

func (s *Server) ListenChans() {
	go func() {
		for {
			select {
			case message_in := <-s.Incoming:
				s.infoLog.Printf("Received MESSAGE: (%s) from CLIENT: (%s); in ROOM: (%s)", strings.Trim(message_in.Text, "\n"), message_in.Author.Username, message_in.Dest.Name)
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

	return fmt.Sprintf("%v %v: %v", msg.Time, msg.Author, msg.Text)

}

func (s *Server) ParseCommand(msg *Message) {
	switch {
	case GetCommand(msg.Text) == 1:
		s.infoLog.Println("Received JOIN command")
		s.JoinRoom(msg.Author, s.GetSecArg(msg))
	case GetCommand(msg.Text) == 2:
		s.infoLog.Println("Received LEAVE command")
		s.LeaveRoom(msg.Author)
	case GetCommand(msg.Text) == 3:
		s.infoLog.Println("Received LIST command")
		s.ListRooms(msg.Author)
	case GetCommand(msg.Text) == 4:
		s.infoLog.Println("Received HELP command")
		s.Help(msg.Author)
	case GetCommand(msg.Text) == 5:
		s.infoLog.Println("Received CREATE command")
		s.CreateRoom(msg.Author, s.GetSecArg(msg))
	}
}

func GetCommand(text string) int {
	command := strings.Split(text, " ")[0]
	command = strings.Trim(command, "\n")
	fmt.Println(command)
	return cmd[command]
}

func (s *Server) GetSecArg(msg *Message) string {
	return strings.Trim(strings.Split(msg.Text, " ")[1], "\n")
}

func (s *Server) LeaveRoom(client *Client) {
	delete(s.Rooms[client.Room.Name].Users, client.Username)
	client.Room = s.Rooms["General"]
	s.Rooms["General"].Users[client.Username] = client
	s.infoLog.Printf("CLIENT: (%s) left ROOM: (%s)", client.Username, client.Room.Name)
}

func (s *Server) Help(client *Client) {
	client.Write("Available commands:\n")
	client.Write("\n")
	client.Write("join <name>\n")
	client.Write("\n")
	client.Write("leave\n")
	client.Write("\n")
	client.Write("list\n")
	client.Write("\n")
	client.Write("create <name>\n")
	client.Write("\n")
}

func (s *Server) JoinRoom(client *Client, name string) {
	if s.Rooms[name] == nil {
		client.Write("ROOM DOES NOT EXIST\n")
		s.warningLog.Printf("CLIENT: (%s) tried to join ROOM: (%s) that does not exist", client.Username, name)
		return
	}

	delete(s.Rooms[client.Room.Name].Users, client.Username)
	client.Room = s.Rooms[name]
	s.Rooms[name].Users[client.Username] = client

	s.infoLog.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room.Name)
}

func (s *Server) CreateRoom(client *Client, name string) {
	if s.Rooms[name] != nil {
		client.Write("ROOM ALREADY EXISTS\n")
		s.warningLog.Printf("CLIENT: (%s) tried to create ROOM: (%s) that already exists", client.Username, client.Room.Name)
		return
	}
	room := NewRoom(name)
	s.Rooms[name] = room
	fmt.Printf("created room: %s", name)
	delete(s.Rooms[client.Room.Name].Users, client.Username)
	client.Room = room
	room.Users[client.Username] = client
	s.infoLog.Printf("CLIENT: (%s) created and joined ROOM: (%s)", client.Username, client.Room.Name)
}

func (s *Server) ListRooms(client *Client) {
	s.infoLog.Printf("CLIENT: (%s) requested list of rooms", client.Username)
	for _, room := range s.Rooms {
		client.Write(room.Name + "\n")
	}
}

func (s *Server) Broadcast(msg *Message) {
	for _, user := range msg.Dest.Users {
		user.Write(s.FormatText(*msg))
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
	Server     *Server
}

func NewClient(connection net.Conn, server *Server) *Client {
	cl := &Client{
		Username:   "anon",
		Reader:     *bufio.NewReader(connection),
		Writer:     *bufio.NewWriter(connection),
		Connection: connection,
		Room:       nil,
		Server:     server,
	}
	server.infoLog.Printf("New CLIENT: %s: %v", cl.Username, connection.RemoteAddr())
	go cl.Read()

	return cl

}

func (client *Client) Read() {

	for {
		str, err := client.Reader.ReadString('\n')
		if err != nil {
			client.Server.errorLog.Println("Reading from CLIENT: ", client.Username)
			return
		}

		msg := NewMessage()
		msg.FillMessage(client, str)
		fmt.Println(client)

		client.Server.Incoming <- msg
	}

}

func (client *Client) Write(str string) {
	_, err := client.Writer.WriteString(str)
	if err != nil {
		return
	}
	client.Server.infoLog.Printf("Sent MESSAGE: (%s) to CLIENT: (%s); in ROOM: (%s)", strings.Trim(str, "\n"), client.Username, client.Room.Name)
	client.Writer.Flush()

}

func (server *Server) Join(client *Client) {
	if server.Rooms["General"].Users[client.Username] != nil || server.Users[client.Username] != nil {
		client.Username = client.Username + "*"
	}
	server.Users[client.Username] = client

	server.Rooms["General"].Users[client.Username] = client
	client.Room = server.Rooms["General"]
	server.infoLog.Printf("CLIENT: (%s) joined ROOM: (%s)", client.Username, client.Room.Name)
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

func main() {

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {

	}

	defer listener.Close()

	server := NewServer()
	server.infoLog.Println("Server started on port 8080")
	room := Room{Name: "General", Users: make(map[string]*Client)}
	server.Rooms[room.Name] = &room
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		nClient := NewClient(conn, server)
		server.Join(nClient)
	}
}
