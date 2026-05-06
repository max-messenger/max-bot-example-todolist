package todorepo

import (
	"github.com/guregu/null/v6"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
)

type todoEntity struct {
	ID        int64
	UserID    null.Int64
	Message   null.String
	Done      null.Bool
	CreatedAt null.Time
}

func (t todoEntity) toDomain() domain.Todo {
	return domain.Todo{
		ID:      t.ID,
		UserID:  t.UserID.Int64,
		Message: t.Message.String,
		Done:    t.Done.Bool,
		Created: t.CreatedAt.Time,
	}
}
