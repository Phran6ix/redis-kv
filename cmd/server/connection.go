package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Phran6ix/redis-kv/internal/commands"
	"github.com/Phran6ix/redis-kv/internal/resp"
)

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	log.Printf("About to start A New Go Routine")

	buf := make([]byte, 1024)

	for {
		fmt.Println("----------------")
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected successfully")
				break
			} else {
				log.Fatal("Could not successfully read data to buffer")
				log.Fatal(err)
			}
			return err
		}

		rawByte := buf[:n]
		commands_args, err := resp.Parse(rawByte)
		if err != nil {
			return err
		}

		done, result, err := commands.ProcessCommand(commands_args)

		var s resp.RedisValue = resp.Null{}
		if err != nil {
			s = resp.Error{Message: err.Error()}
		}

		if done && result != "" {
			switch r := result.(type) {
			case string:
				s = resp.BulkString{Value: r}
			case int:
				s = resp.Integer{Value: r}
			case float64:
				s = resp.Double{Value: r}
			case bool:
				s = resp.Boolean{Value: r}
			}
		}

		connData, err := resp.Serialize(s)
		if err != nil {
			return err
		}

		conn.Write([]byte(connData))
	}
	return nil
}
