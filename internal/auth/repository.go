package auth

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Active       bool
}

type userFinder interface {
	FindByUsername(ctx context.Context, username string) (User, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password, is_active
		FROM `+"`user`"+`
		WHERE email = ?
		LIMIT 1`, username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Active)
	if err != nil {
		return User{}, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}
