package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestTokenIssuerGeneratesValidRS256Token(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	issuer, err := NewTokenIssuer(privatePEM, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	fixedTime := time.Date(2026, 7, 23, 12, 0, 0, 0, time.UTC)
	issuer.now = func() time.Time { return fixedTime }

	rawToken, err := issuer.Generate(User{ID: 2, Username: "User"})
	if err != nil {
		t.Fatal(err)
	}
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (any, error) {
		return &privateKey.PublicKey, nil
	}, jwt.WithValidMethods([]string{"RS256"}), jwt.WithTimeFunc(func() time.Time { return fixedTime }))
	if err != nil || !token.Valid {
		t.Fatalf("parse token: %v", err)
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims["username"] != "User" {
		t.Fatalf("username claim = %v", claims["username"])
	}
	roles, ok := claims["roles"].([]any)
	if !ok || len(roles) != 1 || roles[0] != "ROLE_USER" {
		t.Fatalf("roles claim = %#v", claims["roles"])
	}
}
