package todolistctrl

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router"
)

type listRequest struct {
	UserID int64 `json:"user_id"`
}

type listResponse struct {
	Todos []domain.Todo `json:"todos"`
}

type todoItem struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
	Done    bool   `json:"done"`
	Created string `json:"created"`
}

func (l listResponse) toDomain() []todoItem {
	result := make([]todoItem, 0, len(l.Todos))
	for _, v := range l.Todos {
		result = append(result, todoItem{
			ID:      v.ID,
			Message: v.Message,
			Done:    v.Done,
			Created: v.Created.Format(time.RFC3339),
		})
	}

	return result
}

// listHandler godoc
//
// @Summary		get todo list
// @Description	get all todo records
// @Tags		todo
// @Accept		json
// @Produce		json
// @Success		200		{object} listResponse
// @Failure		400		{object} string
// @Failure		403		{object} string
// @Failure		500		{object} string
// @Router		/v1/todos/list [post]
// .
func (c *Controller) listHandler(w http.ResponseWriter, r *http.Request) {
	var req listRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	records, err := c.srv.List(r.Context(), req.UserID)
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	res := listResponse{
		Todos: records,
	}

	router.WriteSuccess(w, res.toDomain())
}
