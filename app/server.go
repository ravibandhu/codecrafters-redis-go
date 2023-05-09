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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
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
			c.Write([]byte("+PONG\r\n"))
		case "echo":
			c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
		case "set":
			RSet(args[0].String(), args[1].string())
			resp := "ok"
			c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(resp), resp)))
		case "get":
			resp := RGet(args[0].String())
			c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(resp), resp)))
		default:
			c.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
