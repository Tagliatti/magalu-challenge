package httputil

import (
	"encoding/json"
	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
	"net/http"
)

type UnprocessableEntityError struct {
	Errors []string `json:"errors"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func NewErrorMessage(err error) *ErrorMessage {
	return &ErrorMessage{Message: err.Error()}
}

func NewUnprocessableEntityErrorFromZog(errorMessages internals.ZogIssueMap) *UnprocessableEntityError {
	errors := make([]string, 0)

	for path, issues := range errorMessages {
		if path == zconst.ISSUE_KEY_FIRST {
			continue
		}

		for _, issue := range issues {
			errors = append(errors, `The field "`+issue.Path+`" `+issue.Message)
		}
	}

	return &UnprocessableEntityError{
		Errors: errors,
	}
}

func UnprocessableEntityResponse(w http.ResponseWriter, unprocessableEntityError *UnprocessableEntityError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(&unprocessableEntityError)
}

func OkResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&data)
}

func CreatedResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&data)
}

func InternalServerErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&ErrorMessage{err.Error()})
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(&ErrorMessage{err.Error()})
}

func NotFoundResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&ErrorMessage{err.Error()})
}

func NoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
