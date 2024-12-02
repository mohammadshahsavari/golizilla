package middleware

import (
	"fmt"
	"golizilla/config"
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from cookies
		tokenString := c.Cookies("auth_token")
		fmt.Print(tokenString)
		if tokenString == "" {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrMissingAuthToken.Error())
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC and algorithm is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.ErrUnexpectedSigningMethod
			}
			return []byte(cfg.JWTSecretKey), nil
		})

		// Handle token parsing or validation failure
		if err != nil || !token.Valid {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidAuthToken.Error())
		}

		// Extract user ID from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidTokenClaims.Error())
		}

		// Parse user ID from claims
		userID, err := utils.ParseUUID(claims["user_id"].(string))
		if err != nil {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
		}

		// Store user ID in locals for downstream handlers
		c.Locals("user_id", userID)

		// Proceed to the next handler
		return c.Next()
	}
}

func HeaderAuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from cookies
		tokenString := c.Get("auth_token")
		fmt.Print(tokenString)
		if tokenString == "" {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrMissingAuthToken.Error())
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC and algorithm is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.ErrUnexpectedSigningMethod
			}
			return []byte(cfg.JWTSecretKey), nil
		})

		// Handle token parsing or validation failure
		if err != nil || !token.Valid {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidAuthToken.Error())
		}

		// Extract user ID from token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidTokenClaims.Error())
		}

		// Parse user ID from claims
		userID, err := utils.ParseUUID(claims["user_id"].(string))
		if err != nil {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
		}

		// Store user ID in locals for downstream handlers
		c.Locals("user_id", userID)

		// Proceed to the next handler
		return c.Next()
	}
}
