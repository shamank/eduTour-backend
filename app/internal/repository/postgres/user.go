package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/eduTour-backend/app/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	return nil
}

func (r *UserRepo) GetByCredentials(ctx context.Context, email string, passwordHash string) (domain.User, error) {

	return domain.User{}, nil
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	return user, nil
}

func (r *UserRepo) SetRefreshToken(ctx context.Context, userID int, refreshInput domain.RefreshTokenInput) error {
	return nil
}

func (r *UserRepo) Verify(ctx context.Context, userID int) error {
	return nil
}
