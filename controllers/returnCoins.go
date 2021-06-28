package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
	"github.com/priyanshubm/iitk-coin/models"
)

func GetCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getcoins" {
		resp := &models.ServerResponse{
			Message: "Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	rollno, _, _ := jwtTokens.ExtractTokenMetadata(tokenFromUser)
	w.Header().Set("Content-Type", "application/json")

	resp := &models.ServerResponse{
		Message: "",
	}

	switch r.Method {

	case "GET":

		coins, err := jwtTokens.GetCoinsFromRollNo(rollno)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, " -User not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = "Your coins are " + fmt.Sprintf("%f", coins)
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Only GET requests are supported"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
