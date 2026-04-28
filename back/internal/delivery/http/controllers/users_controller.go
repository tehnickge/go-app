package controllers

import (
	"encoding/json"
	"net/http"

	transport "backend-app/internal/delivery/http/transport"
	"backend-app/internal/usecase"

	"backend-app/internal/delivery/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsersController struct {
	uc *usecase.UsersUsecase
}

func NewUsersController(uc *usecase.UsersUsecase) *UsersController {
	return &UsersController{uc: uc}
}

func (c *UsersController) callerIDFromReq(r *http.Request) (uuid.UUID, bool) {
	return middleware.CallerIDFromContext(r.Context())
}

func (c *UsersController) GetMe(w http.ResponseWriter, r *http.Request) {
	callerID, ok := c.callerIDFromReq(r)
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	res, err := c.uc.GetByID(callerID, callerID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

func (c *UsersController) GetByID(w http.ResponseWriter, r *http.Request) {
	callerID, ok := c.callerIDFromReq(r)
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	res, err := c.uc.GetByID(callerID, userID)
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

func (c *UsersController) List(w http.ResponseWriter, r *http.Request) {
	callerID, ok := c.callerIDFromReq(r)
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

type createUserBody struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UsersController) Create(w http.ResponseWriter, r *http.Request) {
	callerID, ok := c.callerIDFromReq(r)
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body createUserBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.uc.Create(callerID, usecase.CreateUserRequest{
		Login:    body.Login,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusCreated, res)
}

type updateUserBody struct {
	Login    *string `json:"login"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Status   *string `json:"status"`
}

func (c *UsersController) Update(w http.ResponseWriter, r *http.Request) {
	callerID, ok := c.callerIDFromReq(r)
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	// Ensure request body is valid JSON (DecodeJSON already does it, but we keep this explicit check).
	var body updateUserBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.uc.Update(callerID, userID, usecase.UpdateUserRequest{
		Login:    body.Login,
		Email:    body.Email,
		Password: body.Password,
		Status:   body.Status,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}

func (c *UsersController) Delete(w http.ResponseWriter, r *http.Request) {
	callerID, ok := c.callerIDFromReq(r)
	if !ok {
		transport.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := c.uc.SoftDelete(callerID, userID); err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

