package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_shortener/configs"
	"url_shortener/internal/db"
	"url_shortener/internal/handlers"
	"url_shortener/internal/logger"
	"url_shortener/internal/repositories"
	"url_shortener/internal/services"
)

func main() {
	errChan := make(chan error, 1)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := configs.LoadConfig()
	l := logger.New(config)
	db := db.ConnectDB(config, l)
	repositories := repositories.New(l, db)
	services := services.New(config, l, repositories)
	services.Start()
	handlers.New(config, l, services)
	server := &http.Server{Addr: config.SERVER_HOST}

	go func() {
		errChan <- server.ListenAndServe()
	}()
	select {
	case err := <-errChan:
		if err != nil {
			l.Error(err)
		}
		services.Stop()

	case s := <-signals:
		l.Info("Received signal", s)
		services.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			l.Info(err)
		}
	}
	l.Info("Server stopped")
}
