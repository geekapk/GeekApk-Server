package main

import (
	"os"
	"net/http"
	"modelmap"
	"models"
	"github.com/jinzhu/gorm"

	// To support databases other than postgres, import the corresponding
	// packages here.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// Open the primary database specified by GA_DB_TYPE and GA_DB_CONN_STR.
	db, err := gorm.Open(
		os.Getenv("GA_DB_TYPE"),
		os.Getenv("GA_DB_CONN_STR"),
	)
	if err != nil {
		panic(err)
	}

	// The secret is supposed to be used for various purposes, e.g. cookie signing.
	secret := os.Getenv("GA_SECRET")
	if len(secret) == 0 {
		panic("GA_SECRET required")
	}

	// A registry that stores models and provides HTTP routing.
	// See src/modelmap/provider.go for a introduction to models.
	registry := modelmap.NewRegistry()

	// Add models here...
	registry.AddProvider(&models.EchoModel {})
	registry.AddProvider(models.NewAccountModel(db))
	registry.AddProvider(models.NewUserModel(db))

	// Build the http.Handler used to serve requests.
	mux, err := registry.BuildHandler(secret)
	if err != nil {
		panic(err)
	}

	// Ignite!
	server := &http.Server {
		Addr: "127.0.0.1:8008",
		Handler: mux,
	}
	server.ListenAndServe()
}
