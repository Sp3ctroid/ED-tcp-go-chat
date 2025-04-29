package model

import (
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	cl_io "clientMod/read_write"
	types "clientMod/types"
)

type createRoomModel struct {
	ta textarea.Model

	conn net.Conn
}

func NewCreateRoomModel(conn *net.Conn) createRoomModel {
	ta := textarea.New()
	ta.Placeholder = "New Room Name..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(1)

	ta.FocusedStyle.Base = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true)
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return createRoomModel{
		ta: ta,

		conn: *conn,
	}
}

func (m createRoomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd tea.Cmd
	m.ta, tiCmd = m.ta.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case types.JSON_payload:
		if strings.Compare("CREATED\n", msg.Status) == 0 {
			m.ta.Placeholder = "New Room Name..."
			return m, tea.Sequence(msgState(types.ChatRoom), JSON_payload_CMD(msg))
		} else if strings.Compare("ALREX", msg.Status) == 0 {
			m.ta.Placeholder = "Room Already Exists"
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			cl_io.WriteToServer(m.conn, "/create "+m.ta.Value()+"\n")
			m.ta.Reset()
			return m, nil
		case tea.KeyEsc:
			return m, msgState(types.CancelJoin)
		}
	}
	return m, tiCmd
}

func (m createRoomModel) View() string {
	return m.ta.View()
}

func (m createRoomModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func JSON_payload_CMD(json_payload types.JSON_payload) tea.Cmd {
	return func() tea.Msg {
		return json_payload
	}
}
