package controller

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Pik-9/fixadated/models"
	"github.com/Pik-9/fixadated/util"
)

func NewEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cevnt models.Event
	err := util.DecodeJsonBody(w, r, &cevnt)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		http.Error(w, "Bad Request", 400)
		return
	}

	ret := models.InsertNewEvent(cevnt)

	w.Header().Set("Content-Type", "application/json")
	w.Write(ret.ToClientJSON(false))
}

func GetEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventid, err := util.Base64ToUuid(ps.ByName("eventid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	event, err := models.GetEventByUuid(eventid)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Event Not Found."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(event.ToClientJSON(true))
}

func PatchEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventid, err := util.Base64ToUuid(ps.ByName("eventid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	event, err := models.GetEventByEditUuid(eventid)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Event not found."))
		return
	}

	var clientData models.Event
	err = util.DecodeJsonBody(w, r, &clientData)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	event.EditSafely(clientData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(event.ToClientJSON(false))
}

func RegisterForEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventid, err := util.Base64ToUuid(ps.ByName("eventid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	event, err := models.GetEventByUuid(eventid)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Event not found."))
		return
	}

	var clientData models.Participant
	err = util.DecodeJsonBody(w, r, &clientData)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	participant, err := event.RegisterParticipant(clientData)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(participant.ToClientJSON(false))
}
