package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func (s *PostgresUserStore) Create(ctx context.Context, user *User) error {
	query := `
INSERT INTO users (username, email, password, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at`
	row := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt)
	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresUserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
SELECT id, username, email, created_at
FROM users
WHERE id = $1`
	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
