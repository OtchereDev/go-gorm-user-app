package middlewares

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
	Secret       string
	Expiry       int64
}

var ConfigDefault = Config{
	Filter:       nil,
	Unauthorized: nil,
	Decode:       nil,
	Secret:       "TEST_SECRET",
	Expiry:       60,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}

	if strings.Trim(cfg.Secret, " ") == "" {
		cfg.Secret = os.Getenv("JWT_SECRET")
	}

	if cfg.Expiry <= 0 {
		cfg.Expiry = ConfigDefault.Expiry
	}

	if cfg.Decode == nil {
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {
			authHeader := c.Get("Authorization")

			if authHeader == "" {
				return nil, errors.New("Authorization headers is required")
			}

			token, err := jwt.Parse(strings.Split(authHeader, " ")[1], func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
				}

				return []byte(cfg.Secret), nil
			})

			if err != nil {
				fmt.Println(err, cfg.Secret)
				return nil, errors.New("Errors parsing token")
			}

			claim, ok := token.Claims.(jwt.MapClaims)

			fmt.Println("token claims: ", claim)

			if !(ok && token.Valid) {
				return nil, errors.New("Invalid Token")
			}

			if expiresAt, ok := claim["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
				return nil, errors.New("jwt is expired")
			}

			return &claim, nil

		}
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}

func New(config Config) fiber.Handler {
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {

		if cfg.Filter != nil && cfg.Filter(c) {
			fmt.Println("AuthMiddleware was skipped")
		}

		fmt.Println("AuthMiddleware was runned")

		claims, err := cfg.Decode(c)

		if err == nil {
			c.Locals("user", *claims)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}
