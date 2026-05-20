package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	template := `<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`
	fmt.Fprintf(w, template, c.fileServerHits.Load())
}

func (c *apiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.fileServerHits.Store(0)
}

func writeError(w http.ResponseWriter, status int, message string) {
	type ErrResponse struct {
		Error string `json:"error"`
	}

	errResponse := ErrResponse{Error: message}
	dat, err := json.Marshal(errResponse)
	if err != nil {
		log.Printf("Error marshalling JSON response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func CleanChirp(chirp string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(chirp, " ")
	for i, word := range words {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				words[i] = "****"
			}
		}
	}
	chirp = strings.Join(words, " ")
	return chirp
}

func (c *apiConfig) HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	if r.Method != "POST" {
		writeError(w, http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Invalid Content-Type")
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(chirp.Body) > 140 {
		writeError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type ValidResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	validResponse := ValidResponse{CleanChirp(chirp.Body)}
	dat, err := json.Marshal(validResponse)
	if err != nil {
		log.Printf("Error marshalling JSON response: %s", err)
		writeError(w, http.StatusInternalServerError, "Error creating JSON response")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func main() {

	serveMux := http.NewServeMux()
	c := &apiConfig{}

	fileServer := http.FileServer(http.Dir("./html"))
	fileServer = http.StripPrefix("/app", fileServer)

	// APP
	serveMux.Handle("/app/", c.IncFileServerHits(fileServer))

	// API
	serveMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	serveMux.HandleFunc("POST /api/validate_chirp", c.HandleValidateChirp)

	serveMux.HandleFunc("GET /api/metrics", c.HandleMetrics)
	serveMux.HandleFunc("POST /api/reset", c.HandleReset)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	fmt.Print("Starting server")
	log.Fatal(server.ListenAndServe())
}
