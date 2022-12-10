package controller

import (
	"net/http"
)

func blankHandler(w *http.ResponseWriter, r http.Request) {
}

func GetRouter() http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("event/edit/", blankHandler)
	mux.HandleFunc("event/register", blankHandler)
	mux.HandleFunc("event/", blankHandler)

	return mux
}
