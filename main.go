package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/controllers"
	"github.com/priyanshubm/iitk-coin/database"
)

func main() {

	http.HandleFunc("/signup", controllers.Handle_signup)
	http.HandleFunc("/login", controllers.Handle_login)
	http.HandleFunc("/secretpage", controllers.Secret_page)
	http.HandleFunc("/addcoins", controllers.AddCoinsHandler)
	http.HandleFunc("/transfercoin", controllers.TransferCoinHandler)
	http.HandleFunc("/getcoins", controllers.GetCoinsHandler)
	log.Printf("Server is up on port 8080")
	err := database.ConnectToDb()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer database.Db.Close()

}
