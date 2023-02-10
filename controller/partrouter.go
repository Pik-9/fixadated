// fixadated is a daemon for a collaborative date finding tool.
//    Copyright (C) 2023 Daniel Steinhauer <d.steinhauer@mailbox.org>
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as
//    published by the Free Software Foundation, either version 3 of the
//    License, or (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.
package controller

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Pik-9/fixadated/models"
	"github.com/Pik-9/fixadated/util"
)

func PatchParticipationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	partid, err := util.Base64ToUuid(ps.ByName("partid"))
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(401)
		w.Write([]byte("Bad Request"))
		return
	}

	participant, err := models.GetParticipantByUuid(partid)
	if err != nil {
		log.Printf("ERROR in %s %s: %s\n", r.Method, r.URL, err.Error())
		w.WriteHeader(404)
		w.Write([]byte("Participant Not Found"))
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

	participant.EditSafely(clientData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(participant.ToClientJSON(false))
}
