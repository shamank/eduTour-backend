package service

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/shamank/eduTour-backend/app/internal/domain"
	"github.com/shamank/eduTour-backend/app/internal/repository"
)

type EventService struct {
	cache *cache.Cache
}

func NewEventService(repo repository.Events, cache *cache.Cache) *EventService {
	return &EventService{
		cache: cache,
	}
}

func (s *EventService) GetAll(ctx context.Context) ([]domain.Event, error) {
	return nil, nil
}
func (s *EventService) Create(ctx context.Context, input EventInput) (int, error) {
	return 0, nil
}

func (s *EventService) GetByID(ctx context.Context, eventID int) (domain.Event, error) {
	return domain.Event{}, nil
}

func (s *EventService) UpdateByID(ctx context.Context, input EventInput) error {
	return nil
}

func (s *EventService) DeleteByID(ctx context.Context, eventID int) error {
	return nil
}
