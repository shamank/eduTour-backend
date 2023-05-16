package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/repository/postgres"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email string, passwordHash string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	SetRefreshToken(ctx context.Context, userID int, refreshInput domain.RefreshTokenInput) error
	Verify(ctx context.Context, userID int) error
}

type Events interface {
	Create(ctx context.Context)
}

type Repository struct {
	db     *sqlx.DB
	Users  Users
	Events Events
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:    db,
		Users: postgres.NewUserRepo(db),
	}
}
