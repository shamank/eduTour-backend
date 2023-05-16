package service

import (
	"context"
	"github.com/shamank/eduTour-backend/app/internal/repository"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
	"github.com/shamank/eduTour-backend/app/pkg/hash"
)

type UserService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewUserService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *UserService {
	return &UserService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *UserService) SignUp(ctx context.Context, input UserSignUpInput) error {
	return nil
}
func (s *UserService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		return Tokens{}, err
	}

	//сделать сохранение рефреш токена в бд
	return Tokens{}, nil
}
func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {

	var tokens Tokens

	return tokens, nil
}

func (s *UserService) Verify(ctx context.Context, userID int, hash string) error {
	return nil
}
