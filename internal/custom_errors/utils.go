package custom_errors

import (
	"encoding/json"
	"log"
	"net/http"
)

func Write400(err string, w http.ResponseWriter) {
	log.Println(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}
