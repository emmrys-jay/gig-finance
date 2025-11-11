package handler

import (
	"errors"
	"net/http"

	"github.com/emmrys-jay/gigmile/internal/service"
	"github.com/emmrys-jay/gigmile/internal/utils"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) GetAccountByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Parse customer ID (handles both GIG prefix and numeric formats)
	id, err := utils.ParseCustomerID(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid customer ID"))
		return
	}

	account, err := h.accountService.GetAccountByCustomer(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

