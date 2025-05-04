package handlers

import (
	"fmt"
	"net/http"
	"time"
	"url_shortener/configs"
	"url_shortener/internal/handlers/url_shortener"
	"url_shortener/internal/logger"
	"url_shortener/internal/middleware"
	"url_shortener/internal/services"

	"github.com/gorilla/mux"
)

func New(config *configs.Config, l logger.Logger, services *services.Service) {
	middleware := middleware.New(config, l, services)
	urlShortenerHandler := url_shortener.New(l, services)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s, %s\n", r.URL.Path, time.Now()) //nolint
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/test_json", urlShortenerHandler.TestJson).Methods("GET")
	api.HandleFunc("/check", urlShortenerHandler.UrlChecker).Methods("GET")
	api.HandleFunc("/short_url", urlShortenerHandler.ShortUrl).Methods("POST")
	api.HandleFunc("/{url:[a-zA-Z0-9]+}", urlShortenerHandler.RedirectUrl).Methods("GET")
	api.Use(middleware.CheckRateLimitMiddleware, middleware.RecoverMiddleware)
	http.Handle("/", r)

}
