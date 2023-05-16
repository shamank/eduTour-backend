package service

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/repository"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
	"github.com/shamank/eduTour-backend/app/pkg/hash"
	"time"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpireIn     time.Duration
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, userID int, hash string) error
}

type EventInput struct {
	Name        string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Categories  []EventCategoryInput
	Cover       string
}

type Events interface {
	GetAll(ctx context.Context) ([]domain.Event, error)
	Create(ctx context.Context, input EventInput) (int, error)
	GetByID(ctx context.Context, eventID int) (domain.Event, error)
	UpdateByID(ctx context.Context, input EventInput) error
	DeleteByID(ctx context.Context, eventID int) error
}

type EventCategoryInput struct {
	Name        string
	Description string
	Slug        string
}

type EventsCategories interface {
	Create(ctx context.Context, input EventCategoryInput) (int, error)
	GetByID(ctx context.Context, categoryID int) (domain.EventCategory, error)
	UpdateByID(ctx context.Context, categoryID int) error
	DeleteByID(ctx context.Context, categoryID int) error
}

type Services struct {
	repos            *repository.Repository
	Users            Users
	Events           Events
	EventsCategories EventsCategories
}

func NewServices(repos *repository.Repository, cache *cache.Cache, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *Services {

	return &Services{
		repos:  repos,
		Users:  NewUserService(repos.Users, hasher, tokenManager),
		Events: NewEventService(repos.Events, cache),
	}
}
