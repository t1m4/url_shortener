package main

import (
	"fmt"
	"log"
	"net/http"
	"url_shortener/configs"
	"url_shortener/internal/db"
	"url_shortener/internal/handlers"
	"url_shortener/internal/repositories"
	"url_shortener/internal/services"
)

func main() {
	config := configs.LoadConfig()
	db := db.ConnectDB(config)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	repositories := repositories.New(db)
	services := services.New(config, repositories)
	handlers.New(services)
	if err := http.ListenAndServe(config.SERVER_HOST, nil); err != nil {
		panic(fmt.Errorf("failed to start server: %s", err))
	}
}
