package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/priyanshubm/iitk-coin/database"
	"github.com/priyanshubm/iitk-coin/models"

	"golang.org/x/crypto/bcrypt"
)

func Handle_signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		resp := &models.ServerResponse{
			Message: "Error:404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

	switch r.Method {

	case "POST":
		var user models.User
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
			resp := &models.ServerResponse{
				Message: "either roll-no or password is empty",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(401)
			resp := &models.ServerResponse{
				Message: "server error",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
		}

		write_err := database.Details(name, rollno, string(hashed_password))

		if write_err != nil {

			w.WriteHeader(500)
			resp := &models.ServerResponse{
				Message: "roll-no is not unique",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := &models.ServerResponse{
			Message: "account created successfully",
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
