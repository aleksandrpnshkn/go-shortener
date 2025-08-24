package handlers

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

func writeBadRequestError(res http.ResponseWriter) {
	writeJSONError(res, http.StatusBadRequest, "bad request")
}

func writeInternalServerError(res http.ResponseWriter) {
	writeJSONError(res, http.StatusInternalServerError, "internal server error")
}

func writeJSONError(res http.ResponseWriter, status int, message string) {
	res.WriteHeader(status)

	responseData := errorResponse{
		Error: apiError{
			Message: message,
		},
	}

	rawResponseData, err := json.Marshal(responseData)
	if err == nil {
		res.Write(rawResponseData)
	}
}
