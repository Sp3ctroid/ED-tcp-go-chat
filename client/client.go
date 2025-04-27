package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
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

		_, err = writer.WriteString(str)
		if err != nil {
			waitGroup.Done()
			return
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
