package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"url_shortener/internal/handlers/url_shortener"
	"url_shortener/internal/services"

	"github.com/gorilla/mux"
)

func New(services *services.Service) {
	urlShortenerHandler := url_shortener.New(services)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s, %s\n", r.URL.Path, time.Now())
		log.Println(time.Now(), r.URL.Query(), r.Method, r.Host)
	})
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()
	r.HandleFunc("/test_json", urlShortenerHandler.TestJson).Methods("GET")
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/check", urlShortenerHandler.UrlChecker).Methods("GET")
	api.HandleFunc("/short_url", urlShortenerHandler.ShortUrl).Methods("POST")
	api.HandleFunc("/{url:[a-zA-Z0-9]+}", urlShortenerHandler.RedirectUrl).Methods("GET")
	http.Handle("/", r)

}
