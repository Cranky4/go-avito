package internalhttp

import (
	"encoding/json"
	"net/http"
)

type DefaultAPIHandler struct{}

func (h *DefaultAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("api documentation comins soon... may be")
}
