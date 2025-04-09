package health

import (
	"net/http"
	"time"
)

type HealthyHandler struct{}

func NewHealthyHandler() *HealthyHandler {
	return &HealthyHandler{}
}
func (h *HealthyHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().UTC().Format(time.RFC3339)))
}
