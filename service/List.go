package service

import (
	"context"
	"errors"

	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
)

type ListService struct {
	repository *repository.ListRepository
}

func NewListService(repository *repository.Repository) *ListService {
	return &ListService{
		repository: repository.List,
	}
}

func (s *ListService) Create(ctx context.Context, listModel model.List) error {
	return s.repository.Create(ctx, listModel)
}

func (s *ListService) Update(ctx context.Context, listModel model.List) (string, error) {
	return s.repository.Update(ctx, listModel)
}

func (s *ListService) Delete(ctx context.Context) (string, error) {
	hasCard, err := s.repository.HasCard(ctx)
	if err != nil {
		return "", err
	}

	if hasCard {
		return "", errors.New("list is not empty")
	}

	return s.repository.Delete(ctx)
}
