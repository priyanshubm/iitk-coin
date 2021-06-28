package controllers

import (
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
)

func Secret_page(w http.ResponseWriter, r *http.Request) {
	resp := &serverResponse{
		Message: "",
	}
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		resp.Message = "Page not formed"
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
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				resp.Message = "Access deined, user not authorized"
				JsonRes, _ := json.Marshal(resp)
				w.Write(JsonRes)

				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		user_roll_no, Acctype, err := jwtTokens.ExtractTokenMetadata(tokenFromUser)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			resp.Message = "Access denied"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		resp.Message = "Welcome to the secret page " + user_roll_no + " " + Acctype
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "only GET requests are supported "
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
	}

}
