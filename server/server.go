package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:4444")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on 127.0.0.1:4444")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a command to execute on the client: ")
		command, _ := reader.ReadString('\n')
		command = command[:len(command)-1] // Remove the newline character
		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}

		fmt.Println("Client response:")
		fmt.Print(response)
	}
}