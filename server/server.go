package main

import (
	"flag"
	"net"
	"serverMod/types"
	"serverMod/types/logger"
	"serverMod/types/rc"
	"serverMod/types/storage"
)

func main() {

	fileLog := flag.Bool("log", false, "=true for logging into file, =false for logging into console")
	ip := flag.String("ip", "127.0.0.1", "ip address. Default is localhost")
	port := flag.String("port", "8080", "port number. Default is 8080")

	flag.Parse()

	listener, err := net.Listen("tcp", *ip+":"+*port)

	if err != nil {

	}

	defer listener.Close()

	server := types.NewServer(*fileLog)
	logger.INFOLOG.Println("Server started on port 8080")
	room := storage.Room{Name: "General", Users: storage.NewClientMap()}
	server.Rooms.CREATE_New_Room(&room)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.ERRORLOG.Println("Couldn't accept connection")
			continue
		}

		nClient := rc.NewClient(conn, server.Incoming)
		server.Join(nClient)
	}
}
