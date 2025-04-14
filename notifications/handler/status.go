package handler

import (
	"errors"
	"github.com/Oudwins/zog"
	"github.com/Tagliatti/magalu-challenge/httputil"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"net/http"
)

var errNotFound = errors.New("notification not found")
var errInvalidOrMissingId = errors.New("invalid or missing notification id")

type StatusHandler struct {
	notificationRepository notifications.Repository
}

func NewStatusHandler(notificationRepository notifications.Repository) *StatusHandler {
	return &StatusHandler{notificationRepository: notificationRepository}
}

func (h *StatusHandler) Handler(w http.ResponseWriter, r *http.Request) {
	var id int64

	validationErrors := zog.Int64().Required().Parse(r.PathValue("id"), &id)

	if validationErrors != nil {
		httputil.BadRequestResponse(w, errInvalidOrMissingId)
		return
	}

	notification, err := h.notificationRepository.FindNotificationStatusByID(r.Context(), id)

	if err != nil {
		httputil.InternalServerErrorResponse(w, err)
		return
	}

	if notification == nil {
		httputil.NotFoundResponse(w, errNotFound)
		return
	}

	httputil.OkResponse(w, &notification)
}
