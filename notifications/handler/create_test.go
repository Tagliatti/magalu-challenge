package handler

import (
	"bytes"
	"encoding/json"
	"github.com/Tagliatti/magalu-challenge/httputil"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"github.com/Tagliatti/magalu-challenge/notifications/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSuccessCreate(t *testing.T) {
	t.Run("Should create a notification successfully", func(t *testing.T) {
		body := bytes.NewBufferString(`{"type":"sms","recipient":"1234567890"}`)
		fixedTime := time.Now().UTC()

		var createNotification notifications.CreateNotification

		err := json.Unmarshal(body.Bytes(), &createNotification)

		require.Nilf(t, err, "Failed to unmarshal JSON: %v", err)

		notification := notifications.Notification{
			Id:        1,
			Type:      createNotification.Type,
			Recipient: createNotification.Recipient,
			CreatedAt: fixedTime,
			Sent:      false,
			SentAt:    nil,
		}

		response := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/notifications", io.NopCloser(body))

		repository := mocks.NewRepository(t)
		repository.On("CreateNotification", &createNotification).Return(int64(1), nil)
		repository.On("FindNotificationByID", int64(1)).Return(&notification, nil)

		NewCreateHandler(repository).
			Handler(response, request)

		expectedStatusCode := http.StatusCreated
		expectedBody, err := json.Marshal(notification)

		require.Nilf(t, err, "Failed to marshal JSON: %v", err)

		assert.Equal(t, expectedStatusCode, response.Code)
		assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
	})
}

func TestInvalidBodyOnCreate(t *testing.T) {
	testCases := []struct {
		name string
		body io.Reader
	}{
		{"Should return 400 when invalid request body (nil body)", nil},
		{"Should return 400 when invalid request body (invalid json)", io.NopCloser(strings.NewReader(`invalid`))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/notifications", tc.body)

			repository := mocks.NewRepository(t)

			NewCreateHandler(repository).
				Handler(response, request)

			expectedStatusCode := http.StatusBadRequest
			expectedBody, err := json.Marshal(httputil.NewErrorMessage(errInvalidBody))

			require.Nilf(t, err, "Failed to marshal JSON: %v", err)

			assert.Equal(t, expectedStatusCode, response.Code)
			assert.Equal(t, string(expectedBody), strings.Trim(response.Body.String(), "\n"))
		})
	}
}

func TestValidationErrorOnCreate(t *testing.T) {
	testCases := []struct {
		name                 string
		expectedBodyContains string
		body                 io.Reader
	}{
		{"Should return 422 when invalid request body (missing type)", `\"type\"`, io.NopCloser(strings.NewReader(`{"recipient":"1234567890"}`))},
		{"Should return 422 when invalid request body (invalid type)", `\"type\"`, io.NopCloser(strings.NewReader(`{"type":"invalid","recipient":"1234567890"}`))},
		{"Should return 422 when invalid request body (missing recipient)", `\"recipient\"`, io.NopCloser(strings.NewReader(`{"type":"sms"}`))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/notifications", tc.body)

			repository := mocks.NewRepository(t)

			NewCreateHandler(repository).
				Handler(response, request)

			assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
			assert.Contains(t, response.Body.String(), tc.expectedBodyContains)
		})
	}
}
