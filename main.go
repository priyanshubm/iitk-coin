package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

type serverResponse struct {
	Message string `json:"message"`
}

func get_hashed_pass(rollno string) string {
	database, _ :=
		sql.Open("sqlite3", "./database/user_details.db")
	rollno_int, _ := strconv.Atoi(rollno)
	sqlStatement := `SELECT password FROM user WHERE rollno= $1;`
	row := database.QueryRow(sqlStatement, rollno_int)

	var hashed_password string
	row.Scan(&hashed_password)
	return (hashed_password)

}
func details(name string, rollno string, password string) error {
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
func Handle_signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		resp := &serverResponse{
			Message: "Error:404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

	switch r.Method {

	case "POST":
		var user User
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		name := user.Name
		rollno := user.Rollno
		password := user.Password
		if rollno == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp := &serverResponse{
				Message: "either roll-no or password is empty",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(401)
			resp := &serverResponse{
				Message: "server error",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
		}

		write_err := details(name, rollno, string(hashed_password))

		if write_err != nil {

			w.WriteHeader(500)
			resp := &serverResponse{
				Message: "roll-no is not unique",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := &serverResponse{
			Message: "account created successfully",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &serverResponse{
			Message: "invalid request, kindly enter a POST request",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
}

func Create_Token(userRollNo string) (string, time.Time, error) {
	var err error
	//creating token

	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("error loading .env file")
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_roll_no"] = userRollNo
	expTime := time.Now().Add(time.Minute * 15)
	atClaims["exp"] = expTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESSKEY")))
	if err != nil {
		return "", time.Now(), err
	}
	return token, expTime, err
}

func Ver_Token(request_token string) (*jwt.Token, error) {
	tokenString := request_token
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("error loading .env file")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESSKEY")), nil //enter secret key
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func Extract_data(user_token string) (string, error) {
	token, err := Ver_Token(user_token)
	if err != nil {
		return " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		roll_no, _ := claims["user_roll_no"].(string)
		return roll_no, err
	}

	return " ", err

}

func Handle_login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		resp := &serverResponse{
			Message: "error:404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {

	case "POST":

		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := user.Rollno
		password := user.Password
		hashedPassword := get_hashed_pass(rollno)
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			w.WriteHeader(500)
			resp := &serverResponse{
				Message: "incorrect password or user not found",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		token, expirationTime, err := Create_Token(rollno)
		if err != nil {
			w.WriteHeader(401)
			resp := &serverResponse{
				Message: "server error",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return

		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)

		resp := &serverResponse{
			Message: "logged in successfully",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &serverResponse{
			Message: "invalid request, kindly enter a POST request",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}

func Secret_page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		resp := &serverResponse{
			Message: "error:404 Page not Found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				resp := &serverResponse{
					Message: "access denied, user not authorized",
				}
				JsonRes, _ := json.Marshal(resp)
				w.Write(JsonRes)

				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		user_roll_no, err := Extract_data(tokenFromUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := &serverResponse{
				Message: "access denied",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		resp := &serverResponse{
			Message: "Welcome" + user_roll_no,
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &serverResponse{
			Message: "invalid request, kindly enter a GET request ",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
	}

}
func main() {

	http.HandleFunc("/signup", Handle_signup)
	http.HandleFunc("/login", Handle_login)
	http.HandleFunc("/secretpage", Secret_page)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
