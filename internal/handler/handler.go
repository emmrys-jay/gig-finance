package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ResponseFormat struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
}

func respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response := ResponseFormat{
		Status:  code >= 200 && code < 300,
		Data:    data,
		Error:   "",
		Message: "operation was successful",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)
}

func respondWithError(w http.ResponseWriter, code int, err error) {

	response := ResponseFormat{
		Status:  false,
		Data:    nil,
		Error:   err.Error(),
		Message: "",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)
}
