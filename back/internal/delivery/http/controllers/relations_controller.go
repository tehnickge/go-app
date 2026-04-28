package controllers

import (
	"net/http"

	transport "backend-app/internal/delivery/http/transport"
	"backend-app/internal/delivery/http/middleware"
	"backend-app/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RelationsController struct {
	uc *usecase.RelationsUsecase
}

func NewRelationsController(uc *usecase.RelationsUsecase) *RelationsController {
	return &RelationsController{uc: uc}
}

func (c *RelationsController) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	res, err := c.uc.GetUserRoles(callerID, userID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

func (c *RelationsController) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	roleID, err := uuid.Parse(chi.URLParam(r, "roleId"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role id"})
		return
	}

	if err := c.uc.AssignRoleToUser(callerID, userID, roleID); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *RelationsController) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}
	roleID, err := uuid.Parse(chi.URLParam(r, "roleId"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role id"})
		return
	}

	if err := c.uc.RemoveRoleFromUser(callerID, userID, roleID); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *RelationsController) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	res, err := c.uc.GetUserGroups(callerID, userID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

type addUserToGroupBody struct {
	Role string `json:"role"`
}

func (c *RelationsController) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid group id"})
		return
	}
	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	var body addUserToGroupBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := c.uc.AddUserToGroup(callerID, groupID, userID, usecase.AddUserToGroupRequest{Role: body.Role}); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *RelationsController) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	callerID, ok := middleware.CallerIDFromContext(r.Context())
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	groupID, err := uuid.Parse(chi.URLParam(r, "groupId"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid group id"})
		return
	}
	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	if err := c.uc.RemoveUserFromGroup(callerID, groupID, userID); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

