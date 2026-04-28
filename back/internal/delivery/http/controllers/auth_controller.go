package controllers

import (
	"net/http"

	transport "backend-app/internal/delivery/http/transport"
	"backend-app/internal/usecase"
)

type AuthController struct {
	auth *usecase.AuthUsecase
}

func NewAuthController(auth *usecase.AuthUsecase) *AuthController {
	return &AuthController{auth: auth}
}

type registerBody struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var body registerBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.auth.Register(usecase.RegisterRequest{
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

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var body loginBody
	if err := transport.DecodeJSON(r, &body); err != nil {
		transport.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := c.auth.Login(usecase.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		transport.WriteUsecaseError(w, err)
		return
	}
	transport.WriteJSON(w, http.StatusOK, res)
}


