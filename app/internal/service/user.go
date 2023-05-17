package service

import (
	"context"
	"github.com/shamank/eduTour-backend/app/internal/domain"
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

	token, err := s.SetRefreshToken(ctx, user.ID)

	//сделать сохранение рефреш токена в бд
	return Tokens{}, nil
}
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (Tokens, error) {

	var tokens Tokens

	return tokens, nil
}

func (s *UserService) Verify(ctx context.Context, userID int, hash string) error {
	return nil
}

func (s *UserService) SetRefreshToken(ctx context.Context, userID int) (string, error) {
	token, expireAt, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	err = s.repo.SetRefreshToken(ctx, userID, domain.RefreshTokenInput{
		RefreshToken: token,
		ExpiresAt:    expireAt,
	})

	return token, err

}
