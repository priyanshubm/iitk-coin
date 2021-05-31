package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	rollno int
	name   string
}

func AddData(db *sql.DB, dt user) {

	st, _ := db.Prepare("INSERT INTO User (rollno, name) VALUES (?, ?)")
	st.Exec(dt.rollno, dt.name)

}

func main() {

	database, _ := sql.Open("sqlite3", "./user_details.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS User ( rollno INTEGER PRIMARY KEY, name TEXT )")
	statement.Exec()

	data := user{rollno: 190655, name: "Priyanshu"}

	AddData(database, data)

}
