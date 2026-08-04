package server

import (
	"io"
	"log"
	"net"

	"github.com/Phran6ix/redis-kv/internal/resp"
)

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	log.Printf("About to start")

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {

		if err == io.EOF {
			log.Fatal("Client disconnected successfully")
		} else {
			log.Fatal("Could not successfully read data to buffer")
			log.Fatal(err)
		}
		return err
	}

	rawByte := buf[:n]
	resp.Parse(rawByte)
	return nil
}
