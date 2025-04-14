package handler

import (
	"encoding/json"
	"errors"
	"github.com/Oudwins/zog"
	"github.com/Tagliatti/magalu-challenge/httputil"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"net/http"
)

var createNotificationSchema = zog.Struct(zog.Schema{
	"type":      zog.String().Trim().Required().OneOf([]string{"email", "sms", "push", "whatsapp"}),
	"recipient": zog.String().Min(3).Max(255).Required(),
})

var errInvalidBody = errors.New("invalid request body")

type CreateHandler struct {
	notificationRepository notifications.Repository
}

func NewCreateHandler(notificationRepository notifications.Repository) *CreateHandler {
	return &CreateHandler{notificationRepository: notificationRepository}
}

func (h *CreateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var createNotification *notifications.CreateNotification
	err := json.NewDecoder(r.Body).Decode(&createNotification)

	if err != nil {
		httputil.BadRequestResponse(w, errInvalidBody)
		return
	}

	validationErrors := createNotificationSchema.Validate(createNotification)

	if validationErrors != nil {
		unprocessableEntityError := httputil.NewUnprocessableEntityErrorFromZog(validationErrors)
		httputil.UnprocessableEntityResponse(w, unprocessableEntityError)
		return
	}

	id, err := h.notificationRepository.CreateNotification(r.Context(), createNotification)

	if err != nil {
		httputil.InternalServerErrorResponse(w, err)
		return
	}

	notification, err := h.notificationRepository.FindNotificationByID(r.Context(), id)

	if err != nil {
		httputil.InternalServerErrorResponse(w, err)
	}

	httputil.CreatedResponse(w, notification)
}
