package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHandlerSuccessfulLogin(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("D7RGR9Sh"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(
		fakeUsers{user: User{ID: 2, Username: "User", PasswordHash: string(hash), Active: true}},
		fakeTokens{token: "jwt-token"},
	)
	handler := NewHandler(service)
	request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"username":"User","password":"D7RGR9Sh"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var body struct {
		Token  string `json:"token"`
		UserID int64  `json:"userId"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Token != "jwt-token" || body.UserID != 2 {
		t.Fatalf("body = %+v", body)
	}
}

func TestHandlerRejectsInvalidBody(t *testing.T) {
	handler := NewHandler(NewService(fakeUsers{}, fakeTokens{}))
	request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"username":"User","unexpected":true}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
