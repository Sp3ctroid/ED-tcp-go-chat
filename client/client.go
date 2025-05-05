package main

import (
	"flag"
	"log"
	"net"

	cl_io "clientMod/read_write"
	model "clientMod/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	ip := flag.String("ip", "127.0.0.1", "default is localhost.")
	port := flag.String("port", "8080", "default is 8080.")

	connection, err := net.Dial("tcp", *ip+":"+*port)
	if err != nil {
		log.Println(err)
	}
	p := tea.NewProgram(model.InitialModel(&connection), tea.WithAltScreen())
	go cl_io.ReadFromServer(connection, p)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
