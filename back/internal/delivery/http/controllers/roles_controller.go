package controllers

import (
	"net/http"

	transport "backend-app/internal/delivery/http/transport"
	"backend-app/internal/delivery/http/middleware"
	"backend-app/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RolesController struct {
	uc *usecase.RolesUsecase
}

func NewRolesController(uc *usecase.RolesUsecase) *RolesController {
	return &RolesController{uc: uc}
}

func (c *RolesController) List(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	res, err := c.uc.List(callerID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

type createRoleBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *RolesController) Create(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body createRoleBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.uc.Create(callerID, usecase.CreateRoleRequest{
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusCreated, res)
}

func (c *RolesController) GetByID(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	res, err := c.uc.GetByID(callerID, roleID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

type updateRoleBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (c *RolesController) Update(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var body updateRoleBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.uc.Update(callerID, roleID, usecase.UpdateRoleRequest{
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

func (c *RolesController) SoftDelete(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := c.uc.SoftDelete(callerID, roleID); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

