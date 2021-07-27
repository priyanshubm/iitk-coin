package controllers

import (
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/jwtTokens"
)

type pendingStruct struct {
	RequestId int
	Rollno    string
	IetmId    int
}

type pendingResponse struct {
	Message string
	List    []pendingStruct
}

func PendingRedeemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/pendingredeems" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	_, Acctype, _ := jwtTokens.ExtractTokenMetadata(tokenFromUser)

	if Acctype != "admin" {
		http.Error(w, "Unauthorized!! admins are allowed ", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var pendingList []pendingStruct

	resp := &pendingResponse{
		Message: "There was some error in fetching the response",
		List:    pendingList,
	}

	switch r.Method {

	case "GET":

		pendingRows, err := jwtTokens.GetPendingRedeems()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for pendingRows.Next() {

			var rowInfo pendingStruct
			err = pendingRows.Scan(&rowInfo.RequestId, &rowInfo.Rollno, &rowInfo.IetmId)
			if err != nil {
				resp.Message = "some error occured "
				resp.List = pendingList
				JsonRes, _ := json.Marshal(resp)
				w.Write(JsonRes)
			}
			pendingList = append(pendingList, rowInfo)
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = "List of pending requests are below. Go to /respondRedeem to respond to pending requests -  "
		resp.List = pendingList
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
