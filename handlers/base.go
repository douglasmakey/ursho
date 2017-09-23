package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

func createResponse(w http.ResponseWriter, r Response) {
	d, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	w.Write(d)
}
