package dbase

import (
	"database/sql"
	"fmt"
	"serverpshh/model"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB

	queryExecMsg string = "INSERT INTO messages (nick, receiver, msg, time) VALUES (?, ?, ?, ?)"
)

func InitDB() {
	var err error

	db, err = sql.Open("sqlite3", "chat.db")
	if err != nil {
		panic(err)
	}

	//Chat
	_, err = db.Exec(queryDB())
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(queryContact())
	if err != nil {
		panic(err)
	}

	fmt.Println("Database created")
}

func AddMsgObj(msg model.ChatMessage) {
	_, err := db.Exec(queryExecMsg, msg.Nick, msg.To, msg.Msg, msg.Timestamp.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
}

func GetMsgObj() []model.ChatMessage {
	rows, err := db.Query("SELECT nick, receiver, msg, time FROM messages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var history []model.ChatMessage

	for rows.Next() {
		var nick, receiver, msg, timeStr string
		rows.Scan(&nick, &receiver, &msg, &timeStr)

		t, _ := time.Parse(time.RFC3339, timeStr)
		history = append(history, model.ChatMessage{
			Nick:      nick,
			To:        receiver,
			Msg:       msg,
			Timestamp: t,
		})
	}
	return history
}

func AddContact(owner string, friend string) {
	query := "INSERT INTO contacts (owner, contact) VALUES (?, ?)"

	_, err := db.Exec(query, owner, friend)
	if err != nil {
		fmt.Println("Error add conts", err)
	}
}

func GetContacts(owner string) []string {
	rows, err := db.Query("SELECT contact FROM contacts WHERE owner=?", owner)
	if err != nil {
		panic(err)
		return nil
	}
	defer rows.Close()

	var contacts []string
	for rows.Next() {
		var friend string
		err := rows.Scan(&friend)
		if err != nil {
			continue
		}

		contacts = append(contacts, friend)
	}
	return contacts
}
