package database

import (
	"strconv"
)

func Get_hashed_password(rollno string) string {

	rollno_int, _ := strconv.Atoi(rollno)
	sqlStatement := `SELECT password FROM user WHERE rollno= $1;`
	row := Db.QueryRow(sqlStatement, rollno_int)

	var hashed_password string
	row.Scan(&hashed_password)

	return (hashed_password)

}
