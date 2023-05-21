package service

import (
	"context"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/repository"
)

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, userName string) (UserProfile, error) {

	res, err := s.repo.GetUserProfile(ctx, userName)
	if err != nil {
		return UserProfile{}, err
	}

	userRoles := make([]string, 0)

	for _, role := range res.Roles {
		userRoles = append(userRoles, role.Name)
	}

	return UserProfile{
		FirstName:  res.FirstName,
		LastName:   res.LastName,
		MiddleName: res.MiddleName,
		Avatar:     res.Avatar,
		Roles:      userRoles,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userName string, user UserProfileInput) error {
	err := s.repo.UpdateUserProfile(ctx, domain.User{
		Username:   userName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Avatar:     user.Avatar,
	})
	if err != nil {
		return err
	}
	return nil
}
