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
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Username:     input.Name,
		Email:        input.Email,
		PasswordHash: passwordHash,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
func (s *UserService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		return Tokens{}, err
	}

	//сделать сохранение рефреш токена в бд

	var userRole = "user"

	for _, role := range user.Roles {
		if role.Name == "admin" {
			userRole = "admin"
		}
	}

	return s.setRefreshToken(ctx, user.ID, userRole)
}
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (Tokens, error) {

	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	var userRole = "user"

	for _, role := range user.Roles {
		if role.Name == "admin" {
			userRole = "admin"
		}
	}

	return s.setRefreshToken(ctx, user.ID, userRole)
}

func (s *UserService) Verify(ctx context.Context, userID int, hash string) error {
	return nil
}

func (s *UserService) setRefreshToken(ctx context.Context, userID int, userRole string) (Tokens, error) {

	accessToken, expireIn, err := s.tokenManager.Generate(userID, userRole)
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, expireAt, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	err = s.repo.SetRefreshToken(ctx, userID, domain.RefreshTokenInput{
		RefreshToken: refreshToken,
		ExpiresAt:    expireAt,
	})

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireIn:     expireIn,
	}, err

}
