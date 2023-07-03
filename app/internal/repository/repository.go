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

type Authorization interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email string, passwordHash string) (domain.User, error)
	GetByUsername(ctx context.Context, username string, passwordHash string) (domain.User, error)

	ConfirmUser(ctx context.Context, confirmToken string) error

	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)

	SetRefreshToken(ctx context.Context, userID int, refreshInput domain.RefreshTokenInput) error
	Verify(ctx context.Context, userID int) error
	GetFullUserInfo(ctx context.Context, userID int) (domain.User, error)
}

type Users interface {
	GetUserProfile(ctx context.Context, userName string) (domain.User, error)
	UpdateUserProfile(ctx context.Context, user domain.User) error
}

type Repository struct {
	db            *sqlx.DB
	Authorization Authorization
	Users         Users
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:            db,
		Authorization: postgres.NewAuthRepo(db),
		Users:         postgres.NewUserRepo(db),
	}
}
