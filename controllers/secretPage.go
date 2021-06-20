package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/priyanshubm/iitk-coin/jwtTokens"
	"github.com/priyanshubm/iitk-coin/models"
)

func Secret_page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		resp := &models.ServerResponse{
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
				resp := &models.ServerResponse{
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
		user_roll_no, err := jwtTokens.Extract_data(tokenFromUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := &models.ServerResponse{
				Message: "access denied",
			}
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		resp := &models.ServerResponse{
			Message: "Welcome" + user_roll_no,
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp := &models.ServerResponse{
			Message: "invalid request, kindly enter a GET request ",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
	}

}
