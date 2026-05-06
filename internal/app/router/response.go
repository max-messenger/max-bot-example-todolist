package router

import (
	"net/http"

	"github.com/max-messenger/max-bot-example-todolist/pkg/marshaler"
)

type response struct {
	Data  any   `json:"data,omitempty"`
	Error error `json:"error,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, r any) {
	writeResponse(w, http.StatusOK, response{Data: r})
}

func WriteError(w http.ResponseWriter, status int, err error) {
	writeResponse(w, status, response{Error: err})
}

func writeResponse(w http.ResponseWriter, status int, r any) {
	b, err := marshaler.MarshalJSONWithoutEscape(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(b); err != nil {
		return
	}
}
