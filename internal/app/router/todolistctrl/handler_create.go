package todolistctrl

import (
	"net/http"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router"
	"github.com/max-messenger/max-bot-example-todolist/pkg/marshaler"
)

type createRequest struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

func (r createRequest) toDomain() domain.Todo {
	return domain.Todo{
		UserID:  r.UserID,
		Message: r.Message,
	}
}

type createResponse struct {
	ID int64 `json:"id"`
}

// createHandler godoc
//
// @Summary		add todo
// @Description	create new todo record
// @Tags		todo
// @Accept		json
// @Produce		json
// @Param       request body     createRequest false "request body"
// @Success		200		{object} createResponse
// @Failure		400		{object} string
// @Failure		403		{object} string
// @Failure		500		{object} string
// @Router		/v1/todos/create [post]
// .
func (c *Controller) createHandler(w http.ResponseWriter, r *http.Request) {
	var request createRequest

	err := marshaler.LoadJSONFromReader(r.Body, &request)
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	id, err := c.srv.Create(r.Context(), request.toDomain())
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	router.WriteSuccess(w, createResponse{ID: id})
}
