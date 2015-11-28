package handlers

import (
	"github.com/gorilla/pat"
)

func GetRouter() *pat.Router {
	r := pat.New()

	r.Get("/health", HealthHandler)
	return r
}
