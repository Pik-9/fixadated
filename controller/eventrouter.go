// fixadated is a daemon for a collaborative date finding tool.
//
//	Copyright (C) 2023 Daniel Steinhauer <d.steinhauer@mailbox.org>
//
//	This program is free software: you can redistribute it and/or modify
//	it under the terms of the GNU Affero General Public License as
//	published by the Free Software Foundation, either version 3 of the
//	License, or (at your option) any later version.
//
//	This program is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even the implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU Affero General Public License for more details.
//
//	You should have received a copy of the GNU Affero General Public License
//	along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
