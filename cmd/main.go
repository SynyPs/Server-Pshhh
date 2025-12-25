package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serverpshh/model"
	"strings"
	"sync"
	"time"

	"serverpshh/dbase"

	"github.com/gorilla/websocket"
)

var (
	clients = make(map[*websocket.Conn]string)
	mu      sync.Mutex

	broadcast = make(chan model.ChatMessage)
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

	mu.Lock()
	clients[conn] = nick
	mu.Unlock()

	fmt.Println(nick + " : Client Connected")

	//Цикл для чтения
	for {
		var msg model.ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Client Disconnected", err)

			mu.Lock()
			delete(clients, conn)
			mu.Unlock()

			conn.Close()
			break
		}
		msg.Nick = nick
		msg.Timestamp = time.Now()
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		dbase.AddMsgObj(msg)
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			continue
		}

		mu.Lock()
		for client, clientNick := range clients {
			if msg.To == "" || clientNick == msg.To || clientNick == msg.Nick {
				err := client.WriteMessage(websocket.TextMessage, b)
				if err != nil {
					log.Println(err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		mu.Unlock()
	}
}
