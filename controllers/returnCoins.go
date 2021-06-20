package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
	"github.com/priyanshubm/iitk-coin/models"
)

func GetCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getcoins" {
		resp := &models.ServerResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	resp := &models.ServerResponse{
		Message: "",
	}

	switch r.Method {

	case "GET":

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query, err := url.ParseQuery(r.URL.RawQuery)
		_, ok := query["rollno"]

		if err != nil || !ok {
			w.WriteHeader(401)
			resp.Message = "Bad request syntax(rollno missing or syntax error )"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		rollno := query["rollno"][0]
		coins, err := jwtTokens.GetCoinsFromRollNo(rollno)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, " -User not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = "Your coins are " + strconv.Itoa(coins)
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Sorry, only GET requests are supported"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
