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
package main

import (
	"context"
	"fmt"
	"github.com/Pik-9/fixadated/controller"
	"github.com/Pik-9/fixadated/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func signalHandler(srv *http.Server, idleConnsClosed chan<- struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, os.Kill)
	<-sigint

	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}

	close(idleConnsClosed)
}

func main() {
	port := 8080
	mux := controller.GetRouter()
	server := http.Server{
		Addr:         fmt.Sprintf("localhost:%d", port),
		Handler:      mux,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 3,
	}

	sigIdle := make(chan struct{})
	go signalHandler(&server, sigIdle)

	log.Printf("Serving on %s...\n", server.Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-sigIdle

	log.Println("Saving to disk...")
	err = models.SaveToDisk()
	if err != nil {
		log.Printf("ERROR while saving: %s\n", err.Error())
	}

	log.Println("Good Byte")
}
