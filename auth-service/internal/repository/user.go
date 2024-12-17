package repository

import (
	"auth-service/internal/models"
	"context"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	query := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.Role)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = $1"
	row := r.db.QueryRowContext(ctx, query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
