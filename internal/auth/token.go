package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenGenerator interface {
	Generate(user User) (string, error)
}

type TokenIssuer struct {
	privateKey *rsa.PrivateKey
	ttl        time.Duration
	now        func() time.Time
}

func NewTokenIssuer(privateKeyPEM []byte, ttl time.Duration) (*TokenIssuer, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse RSA private key: %w", err)
	}
	return &TokenIssuer{privateKey: privateKey, ttl: ttl, now: time.Now}, nil
}

func (i *TokenIssuer) Generate(user User) (string, error) {
	now := i.now().UTC()
	claims := jwt.MapClaims{
		"iat":      now.Unix(),
		"exp":      now.Add(i.ttl).Unix(),
		"roles":    []string{"ROLE_USER"},
		"username": user.Username,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(i.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	return token, nil
}
