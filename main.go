package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter server to connect to: ")
	serverInput, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal("Error reading server")
	}

	serverName := strings.TrimSpace(serverInput)

	if len(serverName) == 0 {
		log.Fatal("Server name cannot be blank")
	}

	fmt.Print("Enter your nickname: ")
	nicknameInput, _ := reader.ReadString('\n')
	fmt.Print("Enter your username: ")
	usernameInput, _ := reader.ReadString('\n')
	fmt.Print("Enter your name: ")
	nameInput, _ := reader.ReadString('\n')

	nickname := strings.TrimSpace(nicknameInput)
	username := strings.TrimSpace(usernameInput)
	realName := strings.TrimSpace(nameInput)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:6667", serverName))
	if err != nil {
		log.Fatal("Error connecting to IRC server. Check your internet connection or make sure the server is correct.")
	}

	defer conn.Close()

	setNicknameCmd := fmt.Sprintf("NICK %s\r\n", nickname)
	_, err = conn.Write([]byte(setNicknameCmd))
	if err != nil {
		log.Fatal(err)
	}
	setNicknameCmd = fmt.Sprintf("USER %s 0 * :%s\r\n", username, realName)
	_, err = conn.Write([]byte(setNicknameCmd))
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		messageFromServer := scanner.Text()

		if strings.HasPrefix(messageFromServer, "PING") {
			pongMessage := strings.Replace(messageFromServer, "PING", "PONG", 1)
			conn.Write([]byte(fmt.Sprintf("%s\r\n", pongMessage)))
		} else {
			fmt.Println(messageFromServer)
		}
	}
}
