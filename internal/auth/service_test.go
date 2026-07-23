package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeUsers struct {
	user User
	err  error
}

func (f fakeUsers) FindByUsername(context.Context, string) (User, error) {
	return f.user, f.err
}

type fakeTokens struct {
	token string
	err   error
}

func (f fakeTokens) Generate(User) (string, error) { return f.token, f.err }

func TestServiceLogin(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("D7RGR9Sh"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		users      fakeUsers
		username   string
		password   string
		wantErr    error
		wantToken  string
		wantUserID int64
	}{
		{
			name:       "valid credentials",
			users:      fakeUsers{user: User{ID: 2, Username: "user@example.com", PasswordHash: string(hash), Active: true}},
			username:   " user@example.com ",
			password:   "D7RGR9Sh",
			wantToken:  "signed-token",
			wantUserID: 2,
		},
		{
			name:     "wrong password",
			users:    fakeUsers{user: User{PasswordHash: string(hash), Active: true}},
			username: "user@example.com",
			password: "wrong",
			wantErr:  ErrInvalidCredentials,
		},
		{
			name:     "inactive user",
			users:    fakeUsers{user: User{PasswordHash: string(hash), Active: false}},
			username: "user@example.com",
			password: "D7RGR9Sh",
			wantErr:  ErrInvalidCredentials,
		},
		{
			name:     "unknown user",
			users:    fakeUsers{err: sql.ErrNoRows},
			username: "unknown@example.com",
			password: "D7RGR9Sh",
			wantErr:  ErrInvalidCredentials,
		},
		{
			name:     "empty username",
			username: " ",
			password: "D7RGR9Sh",
			wantErr:  ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.users, fakeTokens{token: "signed-token"})
			result, err := service.Login(context.Background(), tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Login() error = %v, want %v", err, tt.wantErr)
			}
			if result.Token != tt.wantToken || result.UserID != tt.wantUserID {
				t.Fatalf("Login() result = %+v", result)
			}
		})
	}
}
