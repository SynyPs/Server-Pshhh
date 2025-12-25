package model

import "time"

type ChatMessage struct {
	Type      string    `json:"type"`
	Nick      string    `json:"nick"`
	To        string    `json:"to"`
	Msg       string    `json:"msg"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	Text       string
	SenderName string
}

const (
	TypeMsg           = "msg"
	TypeAddContact    = "add_contact"
	TypeRemoveContact = "remove_contact"
	TypeGetContact    = "get_contact"
	TypeContactList   = "contact_list"
)
