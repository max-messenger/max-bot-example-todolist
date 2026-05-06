package todolist

import (
	"context"
	"fmt"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
)

const messageLimit = 100

//go:generate go tool mockgen -source=todolist.go -destination=todolist_mocks_test.go  -package=todolist_test
type Repository interface {
	Add(ctx context.Context, d domain.Todo) (int64, error)
	List(ctx context.Context, userID int64) ([]domain.Todo, error)
	Delete(ctx context.Context, d domain.Todo) error
	Toggle(ctx context.Context, d domain.Todo) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(ctx context.Context, d domain.Todo) (int64, error) {
	if len([]rune(d.Message)) > messageLimit {
		return 0, fmt.Errorf("text length exceeds maximum allowed. Maximum is 100 characters")
	}
	id, err := s.repo.Add(ctx, d)

	return id, err
}

func (s *Service) Delete(ctx context.Context, d domain.Todo) error {
	err := s.repo.Delete(ctx, d)

	return err
}

func (s *Service) Toggle(ctx context.Context, d domain.Todo) error {
	err := s.repo.Toggle(ctx, d)

	return err
}

func (s *Service) List(ctx context.Context, userID int64) ([]domain.Todo, error) {
	result, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
