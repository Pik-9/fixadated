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
