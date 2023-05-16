package service

import "github.com/shamank/edutour_auth_service/app/internal/repository"

type Services struct {
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{}
}
