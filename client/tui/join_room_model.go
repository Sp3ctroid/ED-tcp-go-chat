package model

import (
	cl_io "clientMod/read_write"
	types "clientMod/types"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type JoinRoomModel struct {
	ta textarea.Model

	conn net.Conn
}

func (jr JoinRoomModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (jr JoinRoomModel) View() string {
	return jr.ta.View()
}

func (jr JoinRoomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd tea.Cmd
	jr.ta, tiCmd = jr.ta.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case types.JSON_payload:
		if strings.Compare("JOINED", msg.Status) == 0 {
			jr.ta.Placeholder = "Choose Room..."
			return jr, tea.Sequence(msgState(types.ChatRoom), JSON_payload_CMD(msg))
		} else if strings.Compare("NEX", msg.Status) == 0 {
			jr.ta.Placeholder = "Room Doesn't Exist"
			return jr, nil
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return jr, tea.Quit
		case tea.KeyEnter:
			cl_io.WriteToServer(jr.conn, "/join "+jr.ta.Value()+"\n")
			jr.ta.Reset()
			return jr, nil
		case tea.KeyEsc:
			return jr, msgState(types.CancelCreate)
		}
	}
	return jr, tiCmd
}

func NewJoinRoomModel(conn *net.Conn) JoinRoomModel {
	ta := textarea.New()
	ta.Placeholder = "Choose Room..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(1)

	ta.FocusedStyle.Base = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true)
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return JoinRoomModel{
		ta: ta,

		conn: *conn,
	}
}
