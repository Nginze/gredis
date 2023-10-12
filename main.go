package main

import (
	"bufio"
	"fmt"
	"strings"

	// "io"
	"net"
)

const (
	STRING  = '+'
	BULK    = '$'
	INTEGER = ':'
	ERROR   = '-'
	ARRAY   = '*'
)

func main() {

	// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

	fmt.Println("Listening on port 6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		// buffer := make([]byte, 1024)
		// _, err := conn.Read(buffer)

		// fmt.Println(string(buffer))

		reader := bufio.NewReader(conn)
		val, err := ParseRESP(reader)

		if err != nil {
			fmt.Println(err)
			return
		}

		command := strings.TrimSpace(strings.ToUpper(val.array[0].bulk))
		args := val.array[1:]

		fmt.Println(command)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Command not found")
			return
		}

		res := handler(args)

		serialized := SerializeRESP(res)

		fmt.Println(string(serialized))

		// fmt.Println(val)

		// if err != nil {
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	fmt.Println(err)
		// 	return
		// }

		conn.Write(serialized)

	}
}
