package cl_io

import (
	"bufio"
	"clientMod/types"
	"encoding/json"
	"net"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

var WaitGroup sync.WaitGroup

func ReadFromServer(connection net.Conn, p *tea.Program) {
	reader := bufio.NewReader(connection)
	json_decoder := json.NewDecoder(reader)
	for {

		payload := types.JSON_payload{}
		json_decoder.Decode(&payload)
		msg := payload
		p.Send(msg)
	}
}

func WriteToServer(connection net.Conn, text string) {
	writer := bufio.NewWriter(connection)

	_, err := writer.WriteString(text)
	if err != nil {
		WaitGroup.Done()
		return
	}

	err = writer.Flush()
	if err != nil {
		WaitGroup.Done()
		return
	}

}
