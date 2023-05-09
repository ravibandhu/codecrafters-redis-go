package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	store := NewStore()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, store)
	}
}

func handleConnection(c net.Conn, store *Store) {
	defer c.Close()
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		value, err := DecodeResp(bufio.NewReader(c))
		if errors.Is(err, io.EOF) {
			break
		}
		fmt.Printf("recevd  %+v\n", &value)
		if err != nil {
			fmt.Println("Error decoding RESP: ", err.Error())
			return
		}
		command := strings.ToLower(value.Array()[0].String())
		args := value.Array()[1:]
		switch command {
		case "ping":
			c.Write(prepareStringResp("PONG"))
		case "echo":
			c.Write(prepareStringRespWithLength(args[0].String()))
		case "set":
			resp, err := store.Set(args)
			if err != nil {
				c.Write(prepareErrResp())
			}
			c.Write(prepareStringResp(resp))
		case "get":
			resp := store.Get(args[0].String())
			if resp != "" {
				c.Write(prepareStringResp(resp))
				return
			}
			c.Write(prepareErrResp())
		default:
			c.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
