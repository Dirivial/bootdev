package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	serveMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./html"))

	serveMux.Handle("/app/", http.StripPrefix("/app", fileServer))

	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	fmt.Print("Starting server")
	log.Fatal(server.ListenAndServe())
}
