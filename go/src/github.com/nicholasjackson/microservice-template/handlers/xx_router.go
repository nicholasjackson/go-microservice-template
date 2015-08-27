package handlers

import (
	"github.com/gorilla/pat"
)

func GetRouter() *pat.Router {
	r := pat.New()

	r.Get("/helloworld", HelloWorldHandler)
	return r
}
