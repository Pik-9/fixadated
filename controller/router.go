package controller

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/Pik-9/fixadated/res"

	"github.com/julienschmidt/httprouter"
)

// TODO: Remove when not needed any longer
func blankHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.Method, r.URL, ps)
	w.WriteHeader(501)
	w.Write([]byte("Not implemented."))
}

func GetRouter() *http.ServeMux {
	appPath, err := fs.Sub(res.Webapp, "webapp")
	if err != nil {
		log.Fatal(err)
	}

	apiMux := httprouter.New()
	apiMux.POST("/api/event", NewEventHandler)
	apiMux.GET("/api/event/:eventid", GetEventHandler)
	apiMux.PATCH("/api/event/:eventid", PatchEventHandler)
	apiMux.POST("/api/event/:eventid/register", RegisterForEventHandler)
	apiMux.PATCH("/api/participation/:partid", PatchParticipationHandler)

	ret := http.NewServeMux()
	ret.Handle("/api/", apiMux)
	ret.Handle("/", http.FileServer(http.FS(appPath)))

	return ret
}
