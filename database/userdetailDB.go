package database

import "database/sql"

func Details(name string, rollno string, password string) error {
	database, _ :=
		sql.Open("sqlite3", "./database/user_details.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno TEXT PRIMARY KEY,password TEXT)")

	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO user (name,rollno,password) VALUES (?, ?, ?)")
	_, err := statement.Exec(name, rollno, password)
	if err != nil {
		return err
	}
	return nil
}
