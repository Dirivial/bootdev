package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (c *apiConfig) IncFileServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d\n", c.fileServerHits.Load())
}

func (c *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.fileServerHits.Store(0)
}

func main() {

	serveMux := http.NewServeMux()
	c := &apiConfig{}

	fileServer := http.FileServer(http.Dir("./html"))
	fileServer = http.StripPrefix("/app", fileServer)

	serveMux.Handle("/app/", c.IncFileServerHits(fileServer))

	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	serveMux.HandleFunc("/metrics", c.HandleMetrics)
	serveMux.HandleFunc("/reset", c.HandleReset)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	fmt.Print("Starting server")
	log.Fatal(server.ListenAndServe())
}
