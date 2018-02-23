package main

import (
	"net/http"
	"modelmap"
	"models"
)

func main() {
	registry := modelmap.NewRegistry()

	// Add models here...
	registry.AddProvider(&models.EchoModel {})

	mux, err := registry.BuildServeMux()
	if err != nil {
		panic(err)
	}

	server := &http.Server {
		Addr: "127.0.0.1:8008",
		Handler: mux,
	}
	server.ListenAndServe()
}
