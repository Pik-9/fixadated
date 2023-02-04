package controller

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/Pik-9/fixadated/res"

	"github.com/julienschmidt/httprouter"
)

func blankHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.Method, r.URL, ps)
	w.WriteHeader(501)
	w.Write([]byte("Not implemented."))
}

func GetRouter() *http.ServeMux {
	webappMux := httprouter.New()
	appPath, err := fs.Sub(res.Webapp, "webapp")
	if err != nil {
		log.Fatal(err)
	}
	webappMux.ServeFiles("/*filepath", http.FS(appPath))

	apiMux := httprouter.New()
	apiMux.POST("/api/event", blankHandler)
	apiMux.GET("/api/event/:eventid", blankHandler)
	apiMux.PATCH("/api/event/:eventid", blankHandler)
	apiMux.POST("/api/event/:eventid/register", blankHandler)
	apiMux.PATCH("/api/participation/:partid", blankHandler)

	ret := http.NewServeMux()
	ret.Handle("/api/", apiMux)
	ret.Handle("/", webappMux)

	return ret
}
