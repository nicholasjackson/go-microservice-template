package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	StatusMessage string
}

func HealthHandler(rw http.ResponseWriter, r *http.Request) {
	response := HealthResponse{}
	response.StatusMessage = "OK"

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}
