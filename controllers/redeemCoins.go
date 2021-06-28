package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
	"github.com/priyanshubm/iitk-coin/models"
)

type redeemCoinsData struct {
	Item_id int `json:"itemid"`
}

func RedeemCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/redeem" {
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

	case "POST":

		var redeemData redeemCoinsData

		err := json.NewDecoder(r.Body).Decode(&redeemData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		item_id := redeemData.Item_id

		if rollno == "" {
			w.WriteHeader(401)
			resp.Message = "Enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		coins, err := jwtTokens.RedeemCoinsDb(rollno, item_id) // withdraw from first user and transfer to second
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp.Message = "Sucessfully redeemed item " + fmt.Sprintf("%d", item_id) + " Coins remaining are " + fmt.Sprintf("%.2f", coins)
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
