package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Print("Starting server")

	serveMux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
