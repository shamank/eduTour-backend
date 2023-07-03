package service

import (
	"context"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/repository"
	"github.com/shamank/eduTour-backend/app/pkg/auth"
	"github.com/shamank/eduTour-backend/app/pkg/email"
	"github.com/shamank/eduTour-backend/app/pkg/hash"
	"net/mail"
)

type AuthService struct {
	repo         repository.Authorization
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	emailManager *email.EmailManager
}

func NewAuthService(repo repository.Authorization, hasher hash.PasswordHasher, tokenManager auth.TokenManager, emailManager *email.EmailManager) *AuthService {
	return &AuthService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
		emailManager: emailManager,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	confirmToken, err := s.tokenManager.GenerateToken(32)
	if err != nil {
		return err
	}

	user := domain.User{
		Username:     input.UserName,
		Email:        input.Email,
		PasswordHash: passwordHash,
		ConfirmToken: confirmToken,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	err = s.emailManager.SendMail([]string{input.Email}, "Password confirm", "Confirm token: "+confirmToken)

	// TODO: сделать обработку ошибки
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	var user domain.User

	_, err = mail.ParseAddress(input.Login)
	if err != nil {
		user, err = s.repo.GetByUsername(ctx, input.Login, passwordHash)
		if err != nil {
			return Tokens{}, err
		}
	} else {
		user, err = s.repo.GetByCredentials(ctx, input.Login, passwordHash)
		if err != nil {
			return Tokens{}, err
		}
	}

	return s.setRefreshToken(ctx, user.ID, user.Username, user.Role.Name)
}

func (s *AuthService) ConfirmUser(ctx context.Context, confirmToken string) error {

	if err := s.repo.ConfirmUser(ctx, confirmToken); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (Tokens, error) {

	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.setRefreshToken(ctx, user.ID, user.Username, user.Role.Name)
}

func (s *AuthService) Verify(ctx context.Context, userID int, hash string) error {
	return nil
}

func (s *AuthService) setRefreshToken(ctx context.Context, userID int, userName string, userRole string) (Tokens, error) {

	accessToken, expireIn, err := s.tokenManager.Generate(userID, userName, userRole)
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
