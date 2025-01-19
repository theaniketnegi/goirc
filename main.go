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

	go func() {
		currChannel := ""
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			// Handle commands like /join, /part, etc.

			if strings.HasPrefix(input, "/join") {
				_, err := conn.Write([]byte(strings.Replace(input, "/join", "JOIN", 1) + "\r\n"))
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				currChannel = strings.Split(input, " ")[1]
			} else if strings.HasPrefix(input, "/part") {
				_, err := conn.Write([]byte(strings.Replace(input, "/part", "PART", 1) + "\r\n"))
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
				currChannel = ""
			} else if strings.HasPrefix(input, "/nick") {
				_, err := conn.Write([]byte(strings.Replace(input, "/nick", "NICK", 1) + "\r\n"))
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
			} else {
				if len(currChannel) == 0 {
					log.Printf("Join a channel before sending a message.\n")
					continue
				}

				_, err := conn.Write([]byte(fmt.Sprintf("PRIVMSG %s :%s", currChannel, input) + "\r\n"))

				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}
		}
	}()

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
