package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Print("Starting server")

	serveMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./html/"))

	serveMux.Handle("/", fileServer)

	assetsServer := http.FileServer(http.Dir("./html/assets"))

	serveMux.Handle("/assets", assetsServer)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
