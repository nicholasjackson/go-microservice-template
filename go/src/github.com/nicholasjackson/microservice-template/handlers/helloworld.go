package handlers

import (
  "net/http"
  "encoding/json"
)

type HelloWorldResponse struct {
	StatusMessage string
}

func HelloWorldHandler(rw http.ResponseWriter, r *http.Request) {
  response := HelloWorldResponse{}
  response.StatusMessage = "Hello World"

  encoder := json.NewEncoder(rw)
	encoder.Encode(&response)
}
