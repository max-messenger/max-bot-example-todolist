package todolistctrl

import (
	"net/http"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router"
	"github.com/max-messenger/max-bot-example-todolist/pkg/marshaler"
)

type deleteRequest struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (r deleteRequest) toDomain() domain.Todo {
	return domain.Todo{
		ID:     r.ID,
		UserID: r.UserID,
	}
}

// deleteHandler godoc
//
// @Summary		delete todo
// @Description	delete by user_id and id
// @Tags		todo
// @Accept		json
// @Produce		json
// @Param       request body     deleteRequest false "request body"
// @Success		200		{object} int64
// @Failure		400		{object} string
// @Failure		403		{object} string
// @Failure		500		{object} string
// @Router		/v1/todos/delete [post]
// .
func (c *Controller) deleteHandler(w http.ResponseWriter, r *http.Request) {
	var request deleteRequest

	err := marshaler.LoadJSONFromReader(r.Body, &request)
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	err = c.srv.Delete(r.Context(), request.toDomain())
	if err != nil {
		router.WriteError(w, http.StatusBadRequest, err)

		return
	}

	router.WriteSuccess(w, request.ID)
}
