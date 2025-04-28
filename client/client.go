package main

import (
	"log"
	"net"

	cl_io "clientMod/read_write"
	model "clientMod/tui"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	move_to_prev_line = "\033[F"
	clear_line        = "\033[K"
)

func main() {
	connection, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println(err)
	}
	p := tea.NewProgram(model.InitialModel(&connection), tea.WithAltScreen())
	go cl_io.ReadFromServer(connection, p)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
