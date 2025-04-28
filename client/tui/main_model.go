package model

import (
	types "clientMod/types"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

func msgState(id types.SessionState) tea.Cmd {
	return func() tea.Msg {
		msg := types.StateChangeMsg{Msg: id}
		return msg
	}
}

type mainModel struct {
	conn  net.Conn
	err   error
	state types.SessionState

	listRoomsModel  tea.Model
	chatRoomModel   tea.Model
	joinRoomModel   tea.Model
	createRoomModel tea.Model
}

func InitialModel(conn *net.Conn) mainModel {

	initialState := types.ChatRoom
	chatRoomModel := NewChatRoomModel(conn)
	createRoomModel := NewCreateRoomModel(conn)
	return mainModel{
		state: initialState,
		conn:  *conn,
		err:   nil,

		createRoomModel: createRoomModel,
		chatRoomModel:   chatRoomModel,
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyCtrlN:
			m.state = types.CreateRoom

		}

	case types.StateChangeMsg:
		switch msg.Msg {
		case types.ChatRoom:
			m.state = msg.Msg
			m.chatRoomModel = NewChatRoomModel(&m.conn)
			m.chatRoomModel, cmd = m.chatRoomModel.Update(tea.WindowSize())
			cmds = append(cmds, cmd)
		case types.CancelCreate:
			m.state = types.ChatRoom
		}

	}

	switch m.state {
	case types.ChatRoom:

		newChatRoom, newCmd := m.chatRoomModel.Update(msg)
		m.chatRoomModel = newChatRoom
		cmd = newCmd

	case types.CreateRoom:

		newCreateRoom, newCmd := m.createRoomModel.Update(msg)
		m.createRoomModel = newCreateRoom
		cmd = newCmd
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.state {
	case types.ChatRoom:
		return m.chatRoomModel.View()
	case types.CreateRoom:
		return m.createRoomModel.View()
	default:
		return m.chatRoomModel.View()
	}

}
