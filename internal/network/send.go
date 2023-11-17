package network

import (
	"log"
	"net"
)

func openSendingServer(port string) net.Listener {
	ln, err := net.Listen("tcp", ":4047")
	if err != nil{
		log.Fatalln("Error opening TCP server:", err)
	}

	return ln
}