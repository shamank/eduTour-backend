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
	UserName string
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

type UserProfile struct {
	UserName   string
	FirstName  string
	LastName   string
	MiddleName string
	Avatar     string
	Roles      []string
}

type Authorization interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, userID int, hash string) error
	setRefreshToken(ctx context.Context, userID int, userName string, userRoles []string) (Tokens, error)
	GetFullUserInfo(ctx context.Context, userID int) (domain.User, error)
}

type UserProfileInput struct {
	FirstName  string
	LastName   string
	MiddleName string
	Avatar     string
}

type Users interface {
	GetUserProfile(ctx context.Context, userName string) (UserProfile, error)
	UpdateUserProfile(ctx context.Context, userName string, user UserProfileInput) error
}

type Services struct {
	repos         *repository.Repository
	Authorization Authorization
	Users         Users
}

func NewServices(repos *repository.Repository, cache *cache.Cache, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *Services {

	return &Services{
		repos:         repos,
		Authorization: NewAuthService(repos.Authorization, hasher, tokenManager),
		Users:         NewUserService(repos.Users),
	}
}
