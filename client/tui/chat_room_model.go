package model

import (
	cl_io "clientMod/read_write"
	types "clientMod/types"
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chatRoomModel struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	conn        net.Conn
	err         error
}

func NewChatRoomModel(conn *net.Conn) chatRoomModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "

	// Remove cursor line styling
	ta.CharLimit = 0
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)

	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).PaddingLeft(3).PaddingTop(1)

	vp.SetContent(`Welcome to General!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chatRoomModel{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true),
		conn:        *conn,
		err:         nil,
	}
}

func (m chatRoomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height("\n\n")
		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case types.JSON_payload:
		if strings.Compare(msg.Status, "BRCREATED") == 0 {
			return m, nil
		}
		m.messages = append(m.messages, msg.Time+" "+m.senderStyle.Render(msg.Username)+" "+msg.Text)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			cl_io.WriteToServer(m.conn, m.textarea.Value()+"\n")
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil
		}

	}
	m.viewport.GotoBottom()

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m chatRoomModel) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		"\n",
		m.textarea.View(),
	)
}

func (m chatRoomModel) Init() tea.Cmd {
	return tea.WindowSize()
}
