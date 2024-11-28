package middleware

import (
	"errors"
	"golizilla/config"
	"golizilla/handler/presenter"
	"golizilla/service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("auth_token")
		if tokenString == "" {
			return presenter.Unauthorized(c, errors.New("missing authentication token"))
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC and algorithm is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			return presenter.Unauthorized(c, errors.New("invalid or expired authentication token"))
		}

		// Extract user ID from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return presenter.Unauthorized(c, errors.New("invalid token claims"))
		}

		userID, err := utils.ParseUUID(claims["user_id"].(string))
		if err != nil {
			return presenter.Unauthorized(c, errors.New("invalid user ID in token"))
		}

		// Store user ID in locals for downstream handlers
		c.Locals("user_id", userID)

		// Proceed to the next handler
		return c.Next()
	}
}
