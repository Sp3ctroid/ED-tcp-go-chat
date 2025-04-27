package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	move_to_prev_line = "\033[F"
	clear_line        = "\033[K"
)

var waitGroup sync.WaitGroup

func ReadFromServer(connection net.Conn) {
	reader := bufio.NewReader(connection)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			waitGroup.Done()
			return
		}

		fmt.Print(str)
	}
}

func WriteToServer(connection net.Conn) {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(os.Stdin)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			waitGroup.Done()
			return
		}

		num, err := writer.WriteString(str)
		if err != nil {
			waitGroup.Done()
			return
		}

		if num != 0 {
			fmt.Print(move_to_prev_line)
			fmt.Print(clear_line)
		}

		err = writer.Flush()
		if err != nil {
			waitGroup.Done()
			return
		}
	}

}

func main() {

	connection, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Println(err)
	}
	waitGroup.Add(1)

	go ReadFromServer(connection)
	go WriteToServer(connection)

	waitGroup.Wait()
}
