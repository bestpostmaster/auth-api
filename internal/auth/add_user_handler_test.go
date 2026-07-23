package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddUserHandlerCreatesUser(t *testing.T) {
	creator := &fakeCreator{userID: 42}
	handler := NewAddUserHandler(NewRegistrationService(creator))
	request := httptest.NewRequest(http.MethodPost, "/api/add-user", bytes.NewBufferString(`{"username":"test","password":"abcd12345"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if strings.TrimSpace(response.Body.String()) != `{"userId":42}` {
		t.Fatalf("body = %s", response.Body.String())
	}
}

func TestAddUserHandlerRejectsDuplicate(t *testing.T) {
	creator := &fakeCreator{err: ErrUsernameAlreadyExists}
	handler := NewAddUserHandler(NewRegistrationService(creator))
	request := httptest.NewRequest(http.MethodPost, "/api/add-user", bytes.NewBufferString(`{"username":"test","password":"abcd12345"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
