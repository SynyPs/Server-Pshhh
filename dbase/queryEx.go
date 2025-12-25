package dbase

func queryDB() string {
	return `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nick Text,
		receiver Text,
		msg Text,
		time Text
	);`
}

func queryContact() string {
	return `
	CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner Text,
		contact Text
	);`
}
