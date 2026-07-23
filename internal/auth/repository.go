package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
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

type userCreator interface {
	Create(ctx context.Context, username, passwordHash, confirmationToken string) (int64, error)
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

func (r *Repository) Create(ctx context.Context, username, passwordHash, confirmationToken string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO `+"`user`"+` (email, password, subscribe_confirmation_token, creation_date)
		SELECT ?, ?, ?, UTC_TIMESTAMP()
		WHERE NOT EXISTS (
			SELECT 1 FROM `+"`user`"+` WHERE email = ?
		)`, username, passwordHash, confirmationToken, username)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			return 0, ErrUsernameAlreadyExists
		}
		return 0, fmt.Errorf("insert user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read inserted rows: %w", err)
	}
	if rowsAffected == 0 {
		return 0, ErrUsernameAlreadyExists
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read inserted user id: %w", err)
	}
	if userID == 0 {
		return 0, errors.New("insert user returned an empty id")
	}
	return userID, nil
}
