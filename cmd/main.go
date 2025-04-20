package main

import (
	"fmt"
	"net/http"
	"url_shortener/configs"
	"url_shortener/internal/db"
	"url_shortener/internal/handlers"
	"url_shortener/internal/logger"
	"url_shortener/internal/repositories"
	"url_shortener/internal/services"
)

func main() {
	config := configs.LoadConfig()
	l := logger.New(config)
	db := db.ConnectDB(config, l)
	repositories := repositories.New(l, db)
	services := services.New(config, l, repositories)
	handlers.New(config, l, services)
	if err := http.ListenAndServe(config.SERVER_HOST, nil); err != nil {
		panic(fmt.Errorf("failed to start server: %s", err))
	}
}
