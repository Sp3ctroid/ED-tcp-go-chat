package model

import (
	"bufio"
	cl_io "clientMod/read_write"
	"clientMod/types"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RoomListModel struct {
	Items  []string
	Cursor int

	height int
	width  int
	conn   net.Conn
}

func (m *RoomListModel) GetAllRooms() {

	cl_io.WriteToServer(m.conn, "/list\n")

	reader := bufio.NewReader(m.conn)
	json_decoder := json.NewDecoder(reader)

	payload := types.JSON_payload{}
	json_decoder.Decode(&payload)
	msg := payload

	m.Items = strings.Split(msg.Text, " ")
}

func NewRoomList(conn *net.Conn) RoomListModel {
	list := RoomListModel{}
	list.Cursor = 0
	list.conn = *conn
	list.GetAllRooms()
	return list
}

func (m RoomListModel) Init() tea.Cmd {
	return nil
}

func (m RoomListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case types.JSON_payload:
		m.Items = append(m.Items, msg.Text)
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}

		case "esc":
			return m, msgState(types.CancelList)

		case "enter":
			room_name := m.Items[m.Cursor]
			cl_io.WriteToServer(m.conn, "/join "+room_name+"\n")
			return m, msgState(types.ChatRoom)
		}

		return m, nil
	}

	return m, nil
}

func (m RoomListModel) View() string {

	s := "Available Rooms\n\n"
	choice_style := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#04B575"))
	cursor_style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575")).Blink(true)
	border_style := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	for i, choice := range m.Items {

		cursor := " "
		if m.Cursor == i {
			cursor = cursor_style.Render("┃ ")
			choice = choice_style.Render(choice)
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, border_style.Render(s))
}
