package main

import (
	"bufio"
	"client/utility"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

var currentChatRoom string
var conn *websocket.Conn
var addr = "ws://localhost:3000/ws"

func main() {
	fmt.Println("Connecting to server...")

	connect()
	var err error
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter your password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	credentials := map[string]string{"username": username, "password": password}
	err = conn.WriteJSON(credentials)
	if err != nil {
		log.Fatal("write:", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
	}
	if string(message) != "Authenticated Successfully" {
		fmt.Println("Authentication failed. Closing connection.")
		return
	}

	fmt.Println("Authenticated successfully.")

	go func() {
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println()
			fmt.Printf("%s: %s\n", msg.From, msg.Message)
		}
	}()

	fmt.Print("Enter message or command (/help for help): ")
	for {
		fmt.Print("Enter message or command : ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "/") && text != "/close" {
			// Handle commands
			command := strings.TrimPrefix(text, "/")
			switch command {
			case "exit":
				fmt.Println("Exiting chat...")
				return
			case "list-chatroom":
				params := make(map[string]string)
				resp, err, _ := utility.CallRestAPIWithMethod(params, "http://"+addr+"/chatroom", "GET")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("List Chatroom Response: " + resp)
			case "list-user":
				params := make(map[string]string)
				resp, err, _ := utility.CallRestAPIWithMethod(params, "http://"+addr+"/user", "GET")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("List User Response: " + resp)
			case "help":
				command := []string{"/help", "/list-chatroom", "/list-user", "/exit", "/join-chatroom", "/leave-chatroom", "/create-chatroom"}
				fmt.Println("All commands", command)
			case "join-chatroom":
				fmt.Print("Chatroom ID : ")
				chatroomID, _ := reader.ReadString('\n')
				chatroomID = strings.TrimSpace(chatroomID)

				params := make(map[string]string)
				params["chatroom_id"] = chatroomID
				params["username"] = username

				resp, err, httpcode := utility.CallRestAPIWithMethod(params, "http://"+addr+"/join-chatroom", "POST")
				if err != nil {
					fmt.Println(err)
				}

				if httpcode != 200 {
					fmt.Println("Failed to join chatroom", resp)
				} else {
					fmt.Println("Join chatroom successfully")
				}
			case "leave-chatroom":
				fmt.Print("Chatroom ID : ")
				chatroomID, _ := reader.ReadString('\n')
				chatroomID = strings.TrimSpace(chatroomID)

				params := make(map[string]string)
				params["chatroom_id"] = chatroomID
				params["username"] = username

				resp, err, httpcode := utility.CallRestAPIWithMethod(params, "http://"+addr+"/leave-chatroom", "POST")
				if err != nil {
					fmt.Println(err)
				}

				if httpcode != 200 {
					fmt.Println("Failed to leave chatroom", resp)
				} else {
					fmt.Println("Leave chatroom successfully")
				}
			case "create-chatroom":
				fmt.Print("Chatroom Name : ")
				chatroomName, _ := reader.ReadString('\n')
				chatroomName = strings.TrimSpace(chatroomName)

				params := make(map[string]string)
				params["chatroom_name"] = chatroomName
				params["username"] = username

				resp, err, httpcode := utility.CallRestAPIWithMethod(params, "http://"+addr+"/chatroom", "POST")
				if err != nil {
					fmt.Println(err)
				}

				if httpcode != 200 {
					fmt.Println("Failed to create chatroom", resp)
				} else {
					fmt.Println("Create chatroom successfully")
				}

			default:
				fmt.Println("Unknown command:", command)
			}
		} else {
			// Send message
			if len(currentChatRoom) == 0 {
				fmt.Print("Ups.. You haven't selected a chatroom/user yet. Please enter username or chatroom : ")
				chatroom, _ := reader.ReadString('\n')
				chatroom = strings.TrimSpace(chatroom)

				currentChatRoom = chatroom
				data := map[string]string{"chatroom": chatroom}
				err = conn.WriteJSON(data)
				if err != nil {
					log.Fatal("write:", err)
				}
			} else {
				msg := Message{
					From:    username,
					To:      currentChatRoom,
					Message: text,
				}
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}
}

func connect() {
	var err error

	conn, _, err = websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
}

func disconnect() {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			log.Println("Error disconnecting:", err)
		} else {
			log.Println("Disconnected from", addr)
		}
	}
}
