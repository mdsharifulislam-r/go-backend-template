package helper

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/mdsharifulislam-r/go-backend-template/config"
	"github.com/mdsharifulislam-r/go-backend-template/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func HashPassword(password string) (string, error) {
	cost := config.Get().BcryptSalt
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func ComparePassword(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func GenerateOTP() int {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		return 1234
	}
	return int(n.Int64()) + 1000
}

func CryptoToken(size int) (string, error) {
	if size <= 0 {
		size = 32
	}
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func AuthToJSON(auth models.Authentication) (datatypes.JSON, error) {
	b, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}

func AuthFromJSON(raw datatypes.JSON) models.Authentication {
	var auth models.Authentication
	if len(raw) == 0 {
		return models.Authentication{IsResetPassword: false}
	}
	_ = json.Unmarshal(raw, &auth)
	return auth
}

func CleanMap(input map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range input {
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		out[k] = v
	}
	return out
}
