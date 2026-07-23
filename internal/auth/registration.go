package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRegistration   = errors.New("invalid username or password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type RegistrationService struct {
	users userCreator
}

func NewRegistrationService(users userCreator) *RegistrationService {
	return &RegistrationService{users: users}
}

func (s *RegistrationService) AddUser(ctx context.Context, username, password string) (int64, error) {
	username = strings.TrimSpace(username)
	if username == "" || len(username) > 250 || len(password) < 8 || len(password) > 72 {
		return 0, ErrInvalidRegistration
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	confirmationToken, err := randomToken()
	if err != nil {
		return 0, err
	}

	userID, err := s.users.Create(ctx, username, string(passwordHash), confirmationToken)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func randomToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate confirmation token: %w", err)
	}
	return hex.EncodeToString(buffer), nil
}
