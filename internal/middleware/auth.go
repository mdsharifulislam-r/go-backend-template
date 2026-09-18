package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mdsharifulislam-r/go-backend-template/internal/errors"
	"github.com/mdsharifulislam-r/go-backend-template/internal/helper"
	"github.com/mdsharifulislam-r/go-backend-template/internal/models"
)

type userKey string

const CurrentUserKey userKey = "currentUser"

func Auth(roles ...models.UserRole) func(http.Handler) http.Handler {
	roleSet := make(map[models.UserRole]struct{}, len(roles))
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				errors.HandleError(w, errors.Unauthorized("You are not authorized"))
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				errors.HandleError(w, errors.Unauthorized("Invalid authorization header"))
				return
			}

			claims, err := helper.ParseToken(parts[1])
			if err != nil {
				errors.HandleError(w, errors.Unauthorized("Invalid or expired token"))
				return
			}

			if len(roleSet) > 0 {
				if _, ok := roleSet[claims.Role]; !ok {
					errors.HandleError(w, errors.Forbidden("You don't have permission"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), CurrentUserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CurrentUser(ctx context.Context) (*helper.Claims, bool) {
	claims, ok := ctx.Value(CurrentUserKey).(*helper.Claims)
	return claims, ok
}
