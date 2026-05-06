package todolist_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/todolist"
)

func TestTodoService(t *testing.T) {
	suite.Run(t, new(testService))
}

type testService struct {
	suite.Suite

	service    *todolist.Service
	repository *MockRepository
}

func (t *testService) SetupTest() {
	ctrl := gomock.NewController(t.T())

	t.repository = NewMockRepository(ctrl)
	t.service = todolist.NewService(t.repository)
}

func (t *testService) TestAddSuccess() {
	ctx := context.Background()
	d := domain.Todo{
		Message: "test",
		Done:    false,
		Created: time.Now(),
	}

	t.repository.EXPECT().Add(ctx, d).Return(int64(1), nil)

	_, err := t.service.Create(ctx, d)
	t.NoError(err)
}

func (t *testService) TestListSuccess() {
	ctx := context.Background()
	d := []domain.Todo{
		{
			Message: "test",
			Done:    false,
			Created: time.Now(),
		},
	}

	t.repository.EXPECT().List(ctx, gomock.Any()).Return(d, nil)

	res, err := t.service.List(ctx, 1)
	t.NoError(err)

	t.Equal(res, d)
}
