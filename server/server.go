package main

import (
	"flag"
	"net"
	"serverMod/types"
)

func main() {

	fileLog := flag.Bool("log", false, "=true for logging into file, =false for logging into console")
	ip := flag.String("ip", "127.0.0.1", "ip address")
	port := flag.String("port", "8080", "port number")

	flag.Parse()

	listener, err := net.Listen("tcp", *ip+":"+*port)

	if err != nil {

	}

	defer listener.Close()

	server := types.NewServer(*fileLog)
	server.InfoLog.Println("Server started on port 8080")
	room := types.Room{Name: "General", Users: make(map[string]*types.Client)}
	server.Rooms[room.Name] = &room
	for {
		conn, err := listener.Accept()
		if err != nil {
			server.ErrorLog.Println("Couldn't accept connection")
			continue
		}

		nClient := types.NewClient(conn, server)
		server.Join(nClient)
	}
}
