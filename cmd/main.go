package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	config, err := configs.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Error while loading config %s", err))
	}
	l := logger.New(config)
	db := db.ConnectDB(config, l)
	repositories := repositories.New(l, db)
	services := services.New(config, l, repositories)
	services.Start()
	handlers.New(config, l, services)
	jsonConfig, _ := json.Marshal(config)
	l.Debug(string(jsonConfig))
	server := &http.Server{
		Addr:         config.App.ServerHost,
		ReadTimeout:  config.App.ReadTimeout,
		WriteTimeout: config.App.WriteTimeout,
		IdleTimeout:  config.App.IdleTimeout,
	}

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
