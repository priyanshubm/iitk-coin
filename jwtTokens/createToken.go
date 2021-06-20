package jwtTokens

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

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
