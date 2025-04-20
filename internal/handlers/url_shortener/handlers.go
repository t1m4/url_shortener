package url_shortener

import (
	"encoding/json"
	"net/http"
	"reflect"
	"url_shortener/internal/custom_errors"
	"url_shortener/internal/logger"
	"url_shortener/internal/schemas"
	"url_shortener/internal/services"

	"github.com/gorilla/mux"
)

type Handler struct {
	l        logger.Logger
	services *services.Service
}

func New(l logger.Logger, services *services.Service) *Handler {
	return &Handler{l, services}
}

func (h *Handler) TestJson(w http.ResponseWriter, r *http.Request) {
	type Message struct {
		Name string
		Body string
		Time int64
	}
	m := Message{"Alice", "Hello", 1294706395881547000}
	b, _ := json.Marshal(m)
	h.l.Debug("Encoded", string(b), b, reflect.TypeOf(b))
	var decodedM Message
	json.Unmarshal(b, &decodedM)
	h.l.Debug("Decoded", decodedM)

	b = []byte(`{"Name":"Bob","Food":"Pickle"}`)
	var decodedString Message
	json.Unmarshal(b, &decodedString)
	h.l.Debug("Decoded string", decodedString)
	json.NewEncoder(w).Encode(m)
}

func (h *Handler) UrlChecker(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	urlStr := query.Get("url")
	if urlStr == "" {
		custom_errors.Write400(h.l, "url is required parameter", w)
		return
	}
	result, err := h.services.URLCheckerService.CheckURL(urlStr)
	if err != nil {
		custom_errors.Write400(h.l, err.Error(), w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)

}

func (h *Handler) ShortUrl(w http.ResponseWriter, r *http.Request) {
	var urlInput schemas.URLInput
	if err := json.NewDecoder(r.Body).Decode(&urlInput); err != nil {
		custom_errors.Write400(h.l, err.Error(), w)
	}
	newUrl, err := h.services.URLShortenerService.ShortURL(&urlInput)
	if err != nil {
		custom_errors.Write400(h.l, err.Error(), w)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"newUrl": newUrl})
}

func (h *Handler) RedirectUrl(w http.ResponseWriter, r *http.Request) {
	newUrl := mux.Vars(r)["url"]
	originalUrl, err := h.services.URLRedirectService.FindRedirectURL(newUrl)
	if err != nil {
		custom_errors.Write400(h.l, err.Error(), w)
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
