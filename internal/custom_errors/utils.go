package custom_errors

import (
	"encoding/json"
	"net/http"
	"url_shortener/internal/logger"
)

func Write400(l logger.Logger, err string, w http.ResponseWriter) {
	l.Debug(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err}) //nolint
}
