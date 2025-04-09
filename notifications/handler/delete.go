package handler

import (
	"github.com/Oudwins/zog"
	"github.com/Tagliatti/magalu-challenge/httputil"
	"github.com/Tagliatti/magalu-challenge/notifications"
	"net/http"
)

type DeleteHandler struct {
	notificationRepository notifications.Repository
}

func NewDeleteHandler(notificationRepository notifications.Repository) *DeleteHandler {
	return &DeleteHandler{notificationRepository: notificationRepository}
}

func (h *DeleteHandler) Handler(w http.ResponseWriter, r *http.Request) {
	var id int64

	validationErrors := zog.Int64().Required().Parse(r.PathValue("id"), &id)

	if validationErrors != nil {
		httputil.BadRequestResponse(w, errInvalidOrMissingId)
		return
	}

	found, err := h.notificationRepository.DeleteNotificationByID(id)

	if err != nil {
		httputil.InternalServerErrorResponse(w, err)
	}

	if !found {
		httputil.NotFoundResponse(w, errNotFound)
		return
	}

	httputil.NoContentResponse(w)
}
