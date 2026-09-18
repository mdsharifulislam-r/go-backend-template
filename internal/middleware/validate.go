package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mdsharifulislam-r/go-backend-template/internal/response"
	"github.com/mdsharifulislam-r/go-backend-template/internal/validation"
)

type bodyKey string

const ValidatedBodyKey bodyKey = "validatedBody"

func Validate[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			raw, err := io.ReadAll(r.Body)
			if err != nil {
				response.Fail(w, http.StatusBadRequest, "Unable to read request body", nil)
				return
			}

			var payload T
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &payload); err != nil {
					response.Fail(w, http.StatusBadRequest, "Invalid JSON body", nil)
					return
				}
			}

			if errs := validation.Struct(payload); errs != nil {
				response.Fail(w, http.StatusBadRequest, "Validation Error", errs)
				return
			}

			ctx := context.WithValue(r.Context(), ValidatedBodyKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BodyFromContext[T any](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(ValidatedBodyKey).(T)
	return val, ok
}
