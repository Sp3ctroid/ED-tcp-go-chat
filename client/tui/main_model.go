package model

import (
	types "clientMod/types"
	"net"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func msgState(id types.SessionState) tea.Cmd {
	return func() tea.Msg {
		msg := types.StateChangeMsg{Msg: id}
		return msg
	}
}

type mainModel struct {
	conn   net.Conn
	err    error
	state  types.SessionState
	height int
	width  int

	listRoomsModel  tea.Model
	chatRoomModel   tea.Model
	joinRoomModel   tea.Model
	createRoomModel tea.Model
	changeNameModel tea.Model
}

func InitialModel(conn *net.Conn) mainModel {

	initialState := types.ChatRoom
	chatRoomModel := NewChatRoomModel(conn)
	createRoomModel := NewCreateRoomModel(conn)
	joinRoomeModel := NewJoinRoomModel(conn)
	listRoomModel := NewRoomList(conn)
	changeNameModel := NewChangeNameModel(conn)
	return mainModel{
		state: initialState,
		conn:  *conn,
		err:   nil,

		listRoomsModel:  listRoomModel,
		joinRoomModel:   joinRoomeModel,
		createRoomModel: createRoomModel,
		chatRoomModel:   chatRoomModel,
		changeNameModel: changeNameModel,
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.chatRoomModel, _ = m.chatRoomModel.Update(msg)
		m.createRoomModel, _ = m.createRoomModel.Update(msg)
		m.joinRoomModel, _ = m.joinRoomModel.Update(msg)
		m.listRoomsModel, _ = m.listRoomsModel.Update(msg)
		m.chatRoomModel, _ = m.chatRoomModel.Update(msg)
		m.changeNameModel, _ = m.changeNameModel.Update(msg)
	case types.JSON_payload:
		if strings.Compare("BRCREATED", msg.Status) == 0 {
			m.listRoomsModel, cmd = m.listRoomsModel.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		//return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyCtrlN:
			m.state = types.CreateRoom
		case tea.KeyCtrlJ:
			m.state = types.JoinRoom
		case tea.KeyCtrlL:
			m.state = types.LeaveRoom
		case tea.KeyCtrlA:
			m.state = types.ListRooms
		case tea.KeyCtrlU:
			m.state = types.ChangeName
		}

	case types.StateChangeMsg:
		switch msg.Msg {
		case types.ChatRoom:
			m.state = msg.Msg
			m.chatRoomModel = NewChatRoomModel(&m.conn)
			newChatRoom, cmd := m.chatRoomModel.Update(tea.WindowSizeMsg{Height: m.height, Width: m.width})
			m.chatRoomModel = newChatRoom
			cmds = append(cmds, cmd)
		case types.CancelCreate:
			m.state = types.ChatRoom
		case types.CancelJoin:
			m.state = types.ChatRoom
		case types.CancelList:
			m.state = types.ChatRoom
		case types.CancelChangeName:
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

	case types.JoinRoom:
		newJoinRoom, newCmd := m.joinRoomModel.Update(msg)
		m.joinRoomModel = newJoinRoom
		cmd = newCmd

	case types.ListRooms:
		newListRoom, newCmd := m.listRoomsModel.Update(msg)
		m.listRoomsModel = newListRoom
		cmd = newCmd

	case types.ChangeName:
		newChangeName, newCmd := m.changeNameModel.Update(msg)
		m.changeNameModel = newChangeName
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
	case types.JoinRoom:
		return m.joinRoomModel.View()
	case types.ListRooms:
		return m.listRoomsModel.View()
	case types.ChangeName:
		return m.changeNameModel.View()
	default:
		return m.chatRoomModel.View()
	}

}
