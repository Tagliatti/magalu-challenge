package handler

import (
	"encoding/json"
	"github.com/Tagliatti/magalu-challenge/httputil"
	"github.com/Tagliatti/magalu-challenge/notifications/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSuccessDelete(t *testing.T) {
	t.Run("Should delete a notification successfully", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/notifications/1", nil)
		request.SetPathValue("id", "1")

		repository := mocks.NewRepository(t)
		repository.On("DeleteNotificationByID", int64(1)).Return(true, nil)

		NewDeleteHandler(repository).
			Handler(response, request)

		expectedStatusCode := http.StatusNoContent

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, "", response.Body.String())
	})
}

func TestNotFoundOnDelete(t *testing.T) {
	t.Run("Should return 404 when notification not found", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/notifications/1", nil)
		request.SetPathValue("id", "1")

		repository := mocks.NewRepository(t)
		repository.On("DeleteNotificationByID", int64(1)).Return(false, nil)

		NewDeleteHandler(repository).
			Handler(response, request)

		expectedStatusCode := http.StatusNotFound
		expectedBody, err := json.Marshal(httputil.NewErrorMessage(errNotFound))

		require.Nilf(t, err, "Failed to marshal JSON: %v", err)

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
	})
}

func TestInvalidIdOnDelete(t *testing.T) {
	t.Run("Should return 400 when id is invalid", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/notifications/invalid-id", nil)
		request.SetPathValue("id", "invalid-id")

		repository := mocks.NewRepository(t)

		NewDeleteHandler(repository).
			Handler(response, request)

		expectedStatusCode := http.StatusBadRequest
		expectedBody, err := json.Marshal(httputil.NewErrorMessage(errInvalidOrMissingId))

		require.Nilf(t, err, "Failed to marshal JSON: %v", err)

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
	})
}
