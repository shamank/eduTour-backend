package service

import (
	"context"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/repository"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
	"github.com/shamank/eduTour-backend/app/pkg/hash"
)

type AuthService struct {
	repo         repository.Authorization
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewAuthService(repo repository.Authorization, hasher hash.PasswordHasher, tokenManager auth.TokenManager) *AuthService {
	return &AuthService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Username:     input.UserName,
		Email:        input.Email,
		PasswordHash: passwordHash,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
func (s *AuthService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		return Tokens{}, err
	}

	//сделать сохранение рефреш токена в бд

	userRoles := make([]string, 0)

	for _, role := range user.Roles {
		userRoles = append(userRoles, role.Name)
	}

	return s.setRefreshToken(ctx, user.ID, user.Username, userRoles)
}
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (Tokens, error) {

	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	userRoles := make([]string, 0)

	for _, role := range user.Roles {
		userRoles = append(userRoles, role.Name)
	}

	return s.setRefreshToken(ctx, user.ID, user.Username, userRoles)
}

func (s *AuthService) Verify(ctx context.Context, userID int, hash string) error {
	return nil
}

func (s *AuthService) setRefreshToken(ctx context.Context, userID int, userName string, userRoles []string) (Tokens, error) {

	accessToken, expireIn, err := s.tokenManager.Generate(userID, userName, userRoles)
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

func (s *AuthService) GetFullUserInfo(ctx context.Context, userID int) (domain.User, error) {
	res, err := s.repo.GetFullUserInfo(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return res, nil
}
