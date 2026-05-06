package todolistctrl

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
)

type TodoService interface {
	Create(ctx context.Context, r domain.Todo) (int64, error)
	List(ctx context.Context, userID int64) ([]domain.Todo, error)
	Delete(ctx context.Context, r domain.Todo) error
	Toggle(ctx context.Context, r domain.Todo) error
}

type Controller struct {
	logger *zap.Logger

	srv TodoService
}

func NewController(l *zap.Logger, s TodoService) *Controller {
	return &Controller{
		logger: l,
		srv:    s,
	}
}

func (c *Controller) Register(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Post("/todos/list", c.listHandler)
		r.Post("/todos/create", c.createHandler)
		r.Post("/todos/delete", c.deleteHandler)
		r.Post("/todos/toggle", c.toggleHandler)
	})
}
