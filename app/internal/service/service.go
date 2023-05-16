package service

import "github.com/shamank/eduTour-backend/app/internal/repository"

type Services struct {
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{}
}
