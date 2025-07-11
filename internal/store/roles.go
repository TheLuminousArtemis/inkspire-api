package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type PostgresRoleStore struct {
	db *sql.DB
}

func (s *PostgresRoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, level, description  FROM roles WHERE name = $1`
	role := &Role{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.Level, &role.Description)
	if err != nil {
		return nil, err
	}
	return role, nil
}
