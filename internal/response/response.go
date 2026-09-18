package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalPage int   `json:"totalPage"`
	Total     int64 `json:"total"`
}

type Body struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message,omitempty"`
	Data       any         `json:"data,omitempty"`
	Errors     any         `json:"errors,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
}

type Options struct {
	Success    bool
	StatusCode int
	Message    string
	Data       any
	Errors     any
	Pagination *Pagination
}

func Send(w http.ResponseWriter, opts Options) {
	if opts.StatusCode == 0 {
		opts.StatusCode = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(opts.StatusCode)
	_ = json.NewEncoder(w).Encode(Body{
		Success:    opts.Success,
		StatusCode: opts.StatusCode,
		Message:    opts.Message,
		Data:       opts.Data,
		Errors:     opts.Errors,
		Pagination: opts.Pagination,
		Timestamp:  time.Now().UTC(),
	})
}

func OK(w http.ResponseWriter, message string, data any) {
	Send(w, Options{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	})
}

func Created(w http.ResponseWriter, message string, data any) {
	Send(w, Options{
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    message,
		Data:       data,
	})
}

func Fail(w http.ResponseWriter, statusCode int, message string, errors any) {
	Send(w, Options{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	})
}
