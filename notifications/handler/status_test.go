package handler

import (
	"encoding/json"
	"github.com/Tagliatti/magalu-challenge/httputil"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"github.com/Tagliatti/magalu-challenge/notifications/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSuccessStatus(t *testing.T) {
	fixedTime := time.Now().UTC()

	testCases := []struct {
		name   string
		status notifications.NotificationStatus
	}{
		{
			"Should return notification status successfully (sent)",
			notifications.NotificationStatus{
				Sent:   false,
				SentAt: nil,
			},
		},
		{
			"Should return notification status successfully (not sent)",
			notifications.NotificationStatus{
				Sent:   true,
				SentAt: &fixedTime,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/notifications/1/status", nil)
			request.SetPathValue("id", "1")

			repository := mocks.NewRepository(t)
			repository.On("FindNotificationStatusByID", int64(1)).Return(&tc.status, nil)

			NewStatusHandler(repository).
				Handler(response, request)

			expectedStatusCode := http.StatusOK
			expectedBody, err := json.Marshal(&tc.status)

			require.Nilf(t, err, "Failed to marshal JSON: %v", err)

			assert.Equal(t, expectedStatusCode, response.Code)
			assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
		})
	}
}

func TestNotFoundOnStatus(t *testing.T) {
	t.Run("Should return 404 when notification not found", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/notifications/1/status", nil)
		request.SetPathValue("id", "1")

		repository := mocks.NewRepository(t)
		repository.On("FindNotificationStatusByID", int64(1)).Return(nil, nil)

		NewStatusHandler(repository).
			Handler(response, request)

		expectedStatusCode := http.StatusNotFound
		expectedBody, err := json.Marshal(httputil.NewErrorMessage(errNotFound))

		require.Nilf(t, err, "Failed to marshal JSON: %v", err)

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
	})
}

func TestInvalidIdOnStatus(t *testing.T) {
	t.Run("Should return 400 when id is invalid", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/notifications/invalid-id/status", nil)
		request.SetPathValue("id", "invalid-id")

		repository := mocks.NewRepository(t)

		NewStatusHandler(repository).
			Handler(response, request)

		expectedStatusCode := http.StatusBadRequest
		expectedBody, err := json.Marshal(httputil.NewErrorMessage(errInvalidOrMissingId))

		require.Nilf(t, err, "Failed to marshal JSON: %v", err)

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
	})
}
