package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/service"
	"github.com/go-playground/validator/v10"
)

type DeploymentHandler struct {
	deploymentService service.DeploymentService
	validator         *validator.Validate
}

func NewDeploymentHandler(deploymentService service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: deploymentService,
		validator:         validator.New(),
	}
}

func (h *DeploymentHandler) RecordDeployment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDeploymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, r, http.StatusBadRequest, err)
		return
	}

	err := h.deploymentService.RecordDeployment(&req)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, r, http.StatusOK, nil)
}
