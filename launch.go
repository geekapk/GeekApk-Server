package main

import (
	"os"
	"net/http"
	"modelmap"
	"models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open(
		os.Getenv("GA_DB_TYPE"),
		os.Getenv("GA_DB_CONN_STR"),
	)
	if err != nil {
		panic(err)
	}

	secret := os.Getenv("GA_SECRET")
	if len(secret) == 0 {
		panic("GA_SECRET required")
	}

	registry := modelmap.NewRegistry()

	// Add models here...
	registry.AddProvider(&models.EchoModel {})
	registry.AddProvider(models.NewAccountModel(db))
	registry.AddProvider(models.NewUserModel(db))

	mux, err := registry.BuildHandler(secret)
	if err != nil {
		panic(err)
	}

	server := &http.Server {
		Addr: "127.0.0.1:8008",
		Handler: mux,
	}
	server.ListenAndServe()
}
