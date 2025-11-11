package handler

import (
	"errors"
	"net/http"

	"github.com/emmrys-jay/gigmile/internal/service"
	"github.com/emmrys-jay/gigmile/internal/utils"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) GetTransactionsByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Parse customer ID (handles both GIG prefix and numeric formats)
	id, err := utils.ParseCustomerID(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid customer ID"))
		return
	}

	transactions, err := h.transactionService.GetTransactionsByCustomer(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, transactions)
}

