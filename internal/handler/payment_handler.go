package handler

import (
	"encoding/json/v2"
	"errors"
	"log"
	"net/http"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/service"
	"github.com/go-playground/validator/v10"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	validator      *validator.Validate
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		validator:      validator.New(),
	}
}

func (h *PaymentHandler) ProcessPaymentNotification(w http.ResponseWriter, r *http.Request) {
	var req models.PaymentNotificationRequest

	if err := json.UnmarshalRead(r.Body, &req); err != nil {
		log.Printf("Error: %v", err)
		respondWithError(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	// Validate request
	err := h.validator.Struct(req)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.paymentService.ProcessPaymentNotification(&req)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, r, http.StatusOK, nil)
}
