package model

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type JoinRoomModel struct {
	ta textarea.Model
}

func (jr JoinRoomModel) Init() tea.Cmd {
	return nil
}

func (jr JoinRoomModel) View() string {
	return jr.View()
}

func (jr JoinRoomModel) Update() (tea.Model, tea.Cmd) {
	return nil, nil
}
