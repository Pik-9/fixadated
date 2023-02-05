package controller

import (
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"

	"github.com/Pik-9/fixadated/util"
	"github.com/Pik-9/fixadated/models"
)

func NewEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

func GetEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventid, err := util.Base64ToUuid(ps.ByName("eventid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	event, err := models.GetEventByUuid(eventid)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Event not found."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(event.GetFlatEvent().ToJSON())
}

func PatchEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventid, err := util.Base64ToUuid(ps.ByName("eventid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	// TODO
	log.Println(r.Method, r.URL, eventid)
	w.WriteHeader(501)
	w.Write([]byte("Not implemented."))
}

func RegisterForEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventid, err := util.Base64ToUuid(ps.ByName("eventid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	// TODO
	log.Println(r.Method, r.URL, eventid)
	w.WriteHeader(501)
	w.Write([]byte("Not implemented."))
}
