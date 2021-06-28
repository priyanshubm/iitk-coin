package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
)

type ItemsData struct {
	Item_id int    `json:"itemid"`
	Cost    string `json:"cost"`
	Number  int    `json:"number"`
}

func AddItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/additems" {
		resp := &serverResponse{
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
	_, Acctype, _ := jwtTokens.ExtractTokenMetadata(tokenFromUser)

	if Acctype == "member" {
		http.Error(w, "Denied, Only Core and admins are allowed ", http.StatusUnauthorized)
		return
	}

	resp := &serverResponse{
		Message: "",
	}

	switch r.Method {

	case "POST":

		var itemData ItemsData

		err := json.NewDecoder(r.Body).Decode(&itemData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item_id := itemData.Item_id

		cost := itemData.Cost
		number := itemData.Number

		w.Header().Set("Content-Type", "application/json")

		message, err := jwtTokens.WriteItemsToDb(item_id, cost, number)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		resp.Message = message + " item added to database "
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Only POST requests are supported"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
