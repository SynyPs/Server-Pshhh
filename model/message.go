package model

import "time"

type ChatMessage struct {
	Nick      string    `json:"nick"`
	To        string    `json:"to"`
	Msg       string    `json:"msg"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	Text       string
	SenderName string
}
