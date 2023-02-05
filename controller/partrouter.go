package controller

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

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

	// TODO
	log.Println(r.Method, r.URL, partid)
	w.WriteHeader(501)
	w.Write([]byte("Not implemented."))
}
