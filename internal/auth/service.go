package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// A fixed valid bcrypt hash keeps the cost of an unknown-user attempt close to
// that of a known-user attempt and makes account enumeration harder.
const dummyPasswordHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoO5uK8eNNoGONhUpwdz3IhG1nHvDaS2nW"

type LoginResult struct {
	Token  string
	UserID int64
}

type Service struct {
	users  userFinder
	tokens tokenGenerator
}

func NewService(users userFinder, tokens tokenGenerator) *Service {
	return &Service{users: users, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResult, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(password))
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}
	passwordMatches := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
	if !user.Active || !passwordMatches {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return LoginResult{}, fmt.Errorf("generate token: %w", err)
	}
	return LoginResult{Token: token, UserID: user.ID}, nil
}
