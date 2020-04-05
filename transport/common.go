package transport

import (
	"encoding/json"
	"net/http"
)

func internalServerError(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": errorMessage,
	})
}
