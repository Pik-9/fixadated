package controller

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func blankHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.Method, r.URL, ps)
	w.WriteHeader(501)
	w.Write([]byte("Not implemented."))
}

func GetRouter() *httprouter.Router {
	mux := httprouter.New()

	mux.POST("/event", blankHandler)
	mux.GET("/event/:eventid", blankHandler)
	mux.PATCH("/event/:eventid", blankHandler)
	mux.POST("/event/:eventid/register", blankHandler)
	mux.PATCH("/participation/:partid", blankHandler)

	return mux
}
