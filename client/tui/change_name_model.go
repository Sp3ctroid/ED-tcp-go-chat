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

type ChangeNameModel struct {
	ta     textarea.Model
	width  int
	height int
	conn   net.Conn
}

func (m ChangeNameModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m ChangeNameModel) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.ta.View())
}

func (m ChangeNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case types.JSON_payload:
		if strings.Compare("CHANGED", msg.Status) == 0 {
			m.ta.Placeholder = "New Name..."
			return m, tea.Sequence(msgState(types.ChatRoom), JSON_payload_CMD(msg))
		} else if strings.Compare("ALRTAK", msg.Status) == 0 {
			m.ta.Placeholder = "Name Already Taken"
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			cl_io.WriteToServer(m.conn, "/name "+m.ta.Value()+"\n")
			m.ta.Reset()
			return m, nil
		case tea.KeyEsc:
			m.ta.Placeholder = "New Name..."
			return m, msgState(types.CancelChangeName)
		}
	}
	m.ta, tiCmd = m.ta.Update(msg)
	return m, tiCmd
}

func NewChangeNameModel(conn *net.Conn) ChangeNameModel {
	ta := textarea.New()
	ta.Placeholder = "New Name..."
	ta.Focus()

	prompt_style := lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Blink(true)
	ta.Prompt = prompt_style.Render("┃ ")
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(1)

	ta.FocusedStyle.Base = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true)
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChangeNameModel{
		ta: ta,

		conn: *conn,
	}
}
