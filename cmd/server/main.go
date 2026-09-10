package server

import (
	"fmt"
	"log"
	"net"
)

func Listener(port string) (net.Listener, error) {
	var address string
	address = "6380"

	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func Start(listener net.Listener) error {
	defer listener.Close()

	for {
		fmt.Println("we are now accepting connections")
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return err
		}

		go handleConnection(conn)
	}
}
