package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/OtchereDev/go-gorm-user-app/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
}

func ComparePassword(h string, p string) error {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
}

func GenerateJWT(u models.UserModel) (string, error) {
	userClaims := jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"email": u.Email,
		"name":  u.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", errors.New("error generating users jwt access")
	}

	return t, nil
}
