package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessRequest(t *testing.T) {
	t.Run("Should return success response", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", nil)

		health := NewHealthyHandler()
		health.Handler(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
		}
		if response.Body.String() == "" {
			t.Errorf("Expected response body, got empty")
		}
	})
}
