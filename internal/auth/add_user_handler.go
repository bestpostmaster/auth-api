package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type AddUserHandler struct {
	service *RegistrationService
}

func NewAddUserHandler(service *RegistrationService) *AddUserHandler {
	return &AddUserHandler{service: service}
}

func (h *AddUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxLoginBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	userID, err := h.service.AddUser(r.Context(), request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidRegistration):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username is required and password must contain 8 to 72 characters"})
		case errors.Is(err, ErrUsernameAlreadyExists):
			writeJSON(w, http.StatusConflict, map[string]string{"error": "username already exists"})
		default:
			log.Printf("add user failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, struct {
		UserID int64 `json:"userId"`
	}{UserID: userID})
}
