package errors

import (
	"errors"
	"log"
	"net/http"

	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func Handle(handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			HandleError(w, err)
		}
	}
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		log.Printf("STATUS: %d | MESSAGE: %s", appErr.StatusCode, appErr.Message)
		response.Fail(w, appErr.StatusCode, appErr.Message, nil)
		return
	}

	log.Printf("UNHANDLED ERROR: %v", err)
	response.Fail(w, http.StatusInternalServerError, "Internal server error", nil)
}
