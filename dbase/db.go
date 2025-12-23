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

	queryExecMsg string = "INSERT INTO messages (nick, msg, time) VALUES (?, ?, ?)"
)

func InitDB() {
	var err error

	db, err = sql.Open("sqlite3", "chat.db")
	if err != nil {
		panic(err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nick Text,
		msg Text,
		time Text
	);`

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database created")
}

func AddMsgObj(msg model.ChatMessage) {
	_, err := db.Exec(queryExecMsg, msg.Nick, msg.Msg, msg.Timestamp.Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
}

func GetMsgObj() []model.ChatMessage {
	rows, err := db.Query("SELECT nick, msg, time FROM messages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var history []model.ChatMessage

	for rows.Next() {
		var nick, msg, timeStr string
		rows.Scan(&nick, &msg, &timeStr)

		t, _ := time.Parse(time.RFC3339, timeStr)
		history = append(history, model.ChatMessage{
			Nick:      nick,
			Msg:       msg,
			Timestamp: t,
		})
	}
	return history
}
