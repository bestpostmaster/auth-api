package auth

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeCreator struct {
	userID            int64
	err               error
	username          string
	passwordHash      string
	confirmationToken string
	calls             int
}

func (f *fakeCreator) Create(_ context.Context, username, passwordHash, confirmationToken string) (int64, error) {
	f.calls++
	f.username = username
	f.passwordHash = passwordHash
	f.confirmationToken = confirmationToken
	return f.userID, f.err
}

func TestRegistrationServiceAddsUser(t *testing.T) {
	creator := &fakeCreator{userID: 42}
	service := NewRegistrationService(creator)

	userID, err := service.AddUser(context.Background(), " test ", "abcd12345")
	if err != nil {
		t.Fatal(err)
	}
	if userID != 42 || creator.username != "test" {
		t.Fatalf("userID = %d, username = %q", userID, creator.username)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(creator.passwordHash), []byte("abcd12345")); err != nil {
		t.Fatalf("stored password is not a valid hash: %v", err)
	}
	if len(creator.confirmationToken) != 64 {
		t.Fatalf("confirmation token length = %d", len(creator.confirmationToken))
	}
}

func TestRegistrationServiceValidatesInput(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{name: "empty username", username: " ", password: "abcd12345"},
		{name: "short password", username: "test", password: "short"},
		{name: "bcrypt password limit", username: "test", password: string(make([]byte, 73))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator := &fakeCreator{}
			_, err := NewRegistrationService(creator).AddUser(context.Background(), tt.username, tt.password)
			if !errors.Is(err, ErrInvalidRegistration) {
				t.Fatalf("error = %v", err)
			}
			if creator.calls != 0 {
				t.Fatalf("Create called %d times", creator.calls)
			}
		})
	}
}

func TestRegistrationServiceReturnsDuplicateError(t *testing.T) {
	creator := &fakeCreator{err: ErrUsernameAlreadyExists}
	_, err := NewRegistrationService(creator).AddUser(context.Background(), "test", "abcd12345")
	if !errors.Is(err, ErrUsernameAlreadyExists) {
		t.Fatalf("error = %v", err)
	}
}
