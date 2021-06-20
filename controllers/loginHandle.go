package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/priyanshubm/iitk-coin/database"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
	"github.com/priyanshubm/iitk-coin/models"

	"golang.org/x/crypto/bcrypt"
)

func Handle_login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		resp := &models.ServerResponse{
			Message: "error:404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {

	case "POST":

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := user.Rollno
		password := user.Password
		hashedPassword := database.Get_hashed_password(rollno)
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			w.WriteHeader(500)
			resp := &models.ServerResponse{
				Message: "incorrect password or user not found",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		token, expirationTime, err := jwtTokens.Create_Token(rollno)
		if err != nil {
			w.WriteHeader(401)
			resp := &models.ServerResponse{
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

		resp := &models.ServerResponse{
			Message: "logged in successfully",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &models.ServerResponse{
			Message: "invalid request, kindly enter a POST request",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
