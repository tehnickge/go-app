package controllers

import (
	"net/http"

	transport "backend-app/internal/delivery/http/transport"
	"backend-app/internal/delivery/http/middleware"
	"backend-app/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GroupsController struct {
	uc *usecase.GroupsUsecase
}

func NewGroupsController(uc *usecase.GroupsUsecase) *GroupsController {
	return &GroupsController{uc: uc}
}

func (c *GroupsController) List(w http.ResponseWriter, r *http.Request) {
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

type createGroupBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *GroupsController) Create(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body createGroupBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.uc.Create(callerID, usecase.CreateGroupRequest{
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusCreated, res)
}

func (c *GroupsController) GetByID(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	res, err := c.uc.GetByID(callerID, groupID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

type updateGroupBody struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (c *GroupsController) Update(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var body updateGroupBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.uc.Update(callerID, groupID, usecase.UpdateGroupRequest{
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

func (c *GroupsController) SoftDelete(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := c.uc.SoftDelete(callerID, groupID); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

