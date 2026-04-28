package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"backend-app/internal/usecase"
)

type errorResponse struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteUsecaseError(w http.ResponseWriter, err error) {
	if err == nil {
		WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "unknown error"})
		return
	}

	switch {
	case errors.Is(err, usecase.ErrUnauthorized):
		WriteJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
	case errors.Is(err, usecase.ErrForbidden):
		WriteJSON(w, http.StatusForbidden, errorResponse{Error: err.Error()})
	case errors.Is(err, usecase.ErrNotFound):
		WriteJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, usecase.ErrConflict):
		WriteJSON(w, http.StatusConflict, errorResponse{Error: err.Error()})
	default:
		// Validation-like errors are currently returned as plain errors.New(...)
		// Best-effort: treat them as 400; everything else becomes 500 only if it looks unexpected.
		// (We avoid leaking internal details.)
		WriteJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	}
}

func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

