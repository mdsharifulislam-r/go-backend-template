package helper

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mdsharifulislam-r/go-backend-template/config"
	"github.com/mdsharifulislam-r/go-backend-template/internal/models"
)

type Claims struct {
	ID    string           `json:"id"`
	Email string           `json:"email"`
	Role  models.UserRole  `json:"role"`
	jwt.RegisteredClaims
}

func SignToken(user *models.User) (string, error) {
	cfg := config.Get()
	expire := parseDuration(cfg.JWTExpire)

	claims := Claims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.Get().JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

func parseDuration(value string) time.Duration {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return 7 * 24 * time.Hour
	}

	if d, err := time.ParseDuration(value); err == nil {
		return d
	}

	if strings.HasSuffix(value, "d") {
		days := strings.TrimSuffix(value, "d")
		var n int
		for _, c := range days {
			if c < '0' || c > '9' {
				return 7 * 24 * time.Hour
			}
			n = n*10 + int(c-'0')
		}
		return time.Duration(n) * 24 * time.Hour
	}

	return 7 * 24 * time.Hour
}
