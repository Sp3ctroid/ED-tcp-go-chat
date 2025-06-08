package rc

import (
	"bufio"
	"encoding/json"
	"net"
	"serverMod/types/json_t"
	"serverMod/types/logger"
	"strings"
	"time"
)

type Client struct {
	Username   string
	Reader     bufio.Reader
	Writer     bufio.Writer
	Connection net.Conn
	Room       string
}

func NewClient(connection net.Conn, message_channel chan *Message) *Client {
	cl := &Client{
		Username:   "anon",
		Reader:     *bufio.NewReader(connection),
		Writer:     *bufio.NewWriter(connection),
		Connection: connection,
		Room:       "",
	}
	logger.INFOLOG.Printf("New CLIENT: %s: %v", cl.Username, connection.RemoteAddr())
	go cl.Read(message_channel)

	return cl

}

func (client *Client) Read(message_channel chan *Message) {

	for {
		str, err := client.Reader.ReadString('\n')
		if err != nil {
			logger.ERRORLOG.Println("Reading from CLIENT: ", client.Username)
			return
		}

		msg := NewMessage()
		msg.FillMessage(client, str)

		message_channel <- msg
	}

}

func (client *Client) Write(Status string, Username string, Text string, Time string) {

	JSON_payload := json_t.JSON_payload{Username: Username, Text: Text, Time: Time, Status: Status}
	str_json, err := json.Marshal(JSON_payload)

	if err != nil {
		return
	}

	_, err = client.Writer.WriteString(string(str_json))
	if err != nil {
		return
	}
	logger.INFOLOG.Printf("GOT PAYLOAD: %s, %s, %s, %s", JSON_payload.Username, JSON_payload.Text, JSON_payload.Time, JSON_payload.Status)
	logger.INFOLOG.Printf("Sent MESSAGE: (%s) to CLIENT: (%s); in ROOM: (%s)", strings.Trim(Text, "\n"), client.Username, client.Room)
	client.Writer.Flush()

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
