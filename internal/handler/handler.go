package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/emmrys-jay/gigmile/internal/middleware"
)

type ResponseFormat struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	start := middleware.GetStartTime(r)

	if data == nil {
		data = struct{}{}
	}

	response := ResponseFormat{
		Status:  code >= 200 && code < 300,
		Data:    data,
		Error:   "",
		Message: "operation was successful",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)

	// Log request
	duration := time.Since(start)
	log.Printf("[%s] %s %s - %d - %v", r.Method, r.URL.Path, r.RemoteAddr, code, duration)
}

func respondWithError(w http.ResponseWriter, r *http.Request, code int, err error) {
	start := middleware.GetStartTime(r)

	response := ResponseFormat{
		Status:  false,
		Data:    struct{}{},
		Error:   err.Error(),
		Message: "",
	}

	log.Printf("Error: %v", err)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("[%s] %s %s - %d - %v", r.Method, r.URL.Path, r.RemoteAddr, http.StatusInternalServerError, time.Since(start))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)

	// Log request
	duration := time.Since(start)
	log.Printf("[%s] %s %s - %d - %v", r.Method, r.URL.Path, r.RemoteAddr, code, duration)
}
