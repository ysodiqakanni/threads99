package helper

import (
	"encoding/json"
	"github.com/ysodiqakanni/threads99/internal/models"
	"net/http"
)

func EncodeSuccessResponse(w http.ResponseWriter, data interface{}, msg string) {
	response := models.NewSuccessResponse(
		data,
		msg,
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func EncodeErrorResponse(w http.ResponseWriter, err error, msg string, errorCode string) {
	response := models.NewErrorOrFailureResponse(
		[]string{err.Error()},
		msg,
		errorCode,
		true,
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// This is to denote unexpected API behavior that needs to be investigated.
func EncodeApiFailureResponse(w http.ResponseWriter, err error, msg string, errorCode string) {
	response := models.NewErrorOrFailureResponse(
		[]string{err.Error()},
		msg,
		errorCode,
		false,
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
