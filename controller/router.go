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
