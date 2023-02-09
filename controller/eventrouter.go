package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/Pik-9/fixadated/models"
	"github.com/Pik-9/fixadated/util"
)

type clientEvent struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Dates       []int64 `json:"dates"`
}

func NewEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cevnt := clientEvent{"", "", make([]int64, 0)}

	err := util.DecodeJsonBody(w, r, &cevnt)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		http.Error(w, "Bad Request", 400)
		return
	}

	dates := make([]time.Time, len(cevnt.Dates))
	for index, tt := range cevnt.Dates {
		dates[index] = time.UnixMilli(tt)
	}

	event := models.NewEvent(cevnt.Name, cevnt.Description, dates)
	eventPlus := models.InsertNewEvent(event)

	w.Header().Set("Content-Type", "application/json")
	w.Write(eventPlus.ToClientJSON())
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

	event, err := models.GetEventByUuid(eventid)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Event not found."))
		return
	}

	var clientData models.Event
	err = util.DecodeJsonBody(w, r, &clientData)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	// The dates cannot be modified afterwards.
	event.Name = clientData.Name
	event.Description = clientData.Description

	w.Header().Set("Content-Type", "application/json")
	w.Write(event.ToClientJSON())
}

func RegisterForEventHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	var clientData models.Participant
	err = util.DecodeJsonBody(w, r, &clientData)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	participant, err := event.RegisterParticipant(clientData)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(participant.ToJSON())
}
