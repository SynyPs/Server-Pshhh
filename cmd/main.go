package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serverpshh/model"
	"strings"
	"time"

	"serverpshh/dbase"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan model.Message)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	dbase.InitDB()
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Listening on port 8443")

	go handleMessages()

	err := http.ListenAndServe("0.0.0.0:8443", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	history := dbase.GetMsgObj()

	for _, msg := range history {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println(err)
			break
		}
	}

	_, b, _ := conn.ReadMessage()
	nick := strings.TrimSpace(string(b))

	clients[conn] = nick
	fmt.Println(nick + " : Client Connected")

	//Цикл для чтения
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client Disconnected", err)
			delete(clients, conn)
			err := conn.Close()
			if err != nil {
				return
			}
			break
		}
		text := string(message)
		broadcast <- model.Message{Text: text, SenderName: nick}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		fmt.Println("Message Received:", msg)

		msgObj := model.ChatMessage{
			Nick:      msg.SenderName,
			Msg:       msg.Text,
			Timestamp: time.Now(),
		}

		dbase.AddMsgObj(msgObj)
		b, _ := json.Marshal(msgObj)

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Printf("Error PUSH msg: %v\n", err)
				err := client.Close()
				if err != nil {
					return
				}
				delete(clients, client)
			}
		}
	}
}
