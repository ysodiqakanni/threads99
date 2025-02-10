package helper

import (
	"encoding/json"
	"github.com/ysodiqakanni/threads99/internal/models"
	"net/http"
)

func EncodeErrorResponse(w http.ResponseWriter, err error, msg string, errorCode string) {
	response := models.NewErrorResponse(
		[]string{err.Error()},
		msg,
		errorCode,
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func EncodeSuccessResponse(w http.ResponseWriter, data interface{}, msg string) {
	response := models.NewSuccessResponse(
		data,
		msg,
	)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
