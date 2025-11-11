package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/service"
	"github.com/emmrys-jay/gigmile/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	customerService service.CustomerService
	validator       *validator.Validate
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		validator:       validator.New(),
	}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customerReq models.CreateCustomerRequest

	if err := json.NewDecoder(r.Body).Decode(&customerReq); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	// Validate request
	if err := h.validator.Struct(customerReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	customer, err := h.customerService.CreateCustomer(&customerReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Parse customer ID (handles both GIG prefix and numeric formats)
	id, err := utils.ParseCustomerID(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid customer ID"))
		return
	}

	customer, err := h.customerService.GetCustomerByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Parse customer ID (handles both GIG prefix and numeric formats)
	id, err := utils.ParseCustomerID(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid customer ID"))
		return
	}

	var customerReq models.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&customerReq); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	// Validate request
	if err := h.validator.Struct(customerReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	customer, err := h.customerService.UpdateCustomer(id, &customerReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Parse customer ID (handles both GIG prefix and numeric formats)
	id, err := utils.ParseCustomerID(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid customer ID"))
		return
	}

	err = h.customerService.DeleteCustomer(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
