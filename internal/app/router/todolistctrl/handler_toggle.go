package todolistctrl

import (
	"net/http"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router"
	"github.com/max-messenger/max-bot-example-todolist/pkg/marshaler"
)

type toggleRequest struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Done   bool  `json:"done"`
}

func (r toggleRequest) toDomain() domain.Todo {
	return domain.Todo{
		ID:     r.ID,
		Done:   r.Done,
		UserID: r.UserID,
	}
}

// toggleHandler godoc
//
// @Summary		toggle todo
// @Tags		todo
// @Accept		json
// @Produce		json
// @Param       request body     toggleRequest false "request body"
// @Success		200		{object} int64
// @Failure		400		{object} string
// @Failure		403		{object} string
// @Failure		500		{object} string
// @Router		/v1/todos/toggle [post]
// .
func (c *Controller) toggleHandler(w http.ResponseWriter, r *http.Request) {
	var request toggleRequest

	err := marshaler.LoadJSONFromReader(r.Body, &request)
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	err = c.srv.Toggle(r.Context(), request.toDomain())
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	router.WriteSuccess(w, request.ID)
}
