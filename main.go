package main

import (
	"fmt"
	"github.com/Pik-9/fixadated/controller"
	"log"
	"net/http"
	"time"
)

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

	log.Printf("Serving on %s...\n", server.Addr)
	err := server.ListenAndServe()
	if err.Error() != "ErrServerClosed" {
		log.Fatal(err)
	}
}
