// server.go

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	banner := `
	   ____________                 
	  / ____/_  __/                 
	 / / __  / /                    
	/ /_/ / / /                     
	\____/ /_/                      
									
	   ____  _____________          
	  / __ \/ __/ __/ ___/___  _____
	 / / / / /_/ /_ \__ \/ _ \/ ___/
	/ /_/ / __/ __/___/ /  __/ /__  
	\____/_/ /_/  /____/\___/\___/  
			`
	red := "\x1b[0;31m"
	reset := "\x1b[0m"

	coloredBanner := red + banner + reset
	fmt.Println(coloredBanner)

	ip := "127.0.0.1:4444"

	listener, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on", ip)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("CMD: ")
		command, _ := reader.ReadString('\n')
		command = command[:len(command)-1] // Remove the newline character

		if command == "exit" {
			fmt.Println("Exiting server.")
			return
		}

		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		readResponse(conn)
	}
}

func readResponse(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}

		response = response[:len(response)-1] // Remove the newline character

		if response == "END_OF_RESPONSE" {
			return
		}

		fmt.Println(response)
	}
}
