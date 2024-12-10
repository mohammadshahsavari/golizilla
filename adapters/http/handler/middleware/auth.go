package middleware

import (
	"golizilla/adapters/http/handler/presenter"
	"golizilla/adapters/persistence/logger"
	"golizilla/config"
	"golizilla/core/service/utils"
	"golizilla/internal/apperrors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var Store *session.Store

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from cookies
		tokenString := c.Cookies("auth_token")

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
		tokenString := c.Get("auth_token")

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
		UserIDJWT, err := utils.ParseUUID(claims["user_id"].(string))
		if err != nil {
			return presenter.SendError(c, fiber.StatusUnauthorized, apperrors.ErrInvalidUserID.Error())
		}

		// // Get the session
		// sess, err := Store.Get(c)
		// if err != nil {
		// 	return presenter.SendError(c, fiber.StatusUnauthorized, "Session not found")
		// }

		// // Retrieve the user ID from the session
		// userIDValue := sess.Get("user_id")
		// if userIDValue == nil {
		// 	return presenter.SendError(c, fiber.StatusUnauthorized, "Unauthorized access")
		// }

		// // Assert the type of userIDValue
		// userID, err := uuid.Parse(userIDValue.(string))
		// if err != nil {
		// 	return presenter.SendError(c, fiber.StatusUnauthorized, "Invalid session data")
		// }

		// // check for CSRF
		// if UserIDJWT != userID {
		// 	return presenter.SendError(c, fiber.StatusUnauthorized, "Unauthorized access")
		// }

		// Store user ID in locals for downstream handlers
		c.Locals("user_id", UserIDJWT)

		// Proceed to the next handler
		return c.Next()
	}
}

func InitSessionStore(cfg *config.Config) {
	Store = session.New(session.Config{
		// Set options for the session store
		CookieHTTPOnly: true,
		CookieSameSite: "Strict",
		Expiration:     cfg.JWTExpiresIn * time.Second,
	})
}

// ContextMiddleware adds trace_id, transaction_id, and user_id to fasthttp.RequestCtx.
func ContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate trace ID and transaction ID
		traceID := c.Get("X-Trace-ID", uuid.New().String())
		transactionID := uuid.New().String()

		// Get user ID from locals
		userID, ok := c.Locals("user_id").(uuid.UUID)
		if !ok {
			userID = uuid.Nil
		}

		sessionID := c.Cookies("session_id", "no-session")
		endpoint := c.OriginalURL()

		// Store values directly in fasthttp.RequestCtx
		c.Context().SetUserValue(logger.TraceIDKey, traceID)
		c.Context().SetUserValue(logger.TransactionIDKey, transactionID)
		c.Context().SetUserValue(logger.UserIDKey, userID.String())
		c.Context().SetUserValue(logger.SessionIDKey, sessionID)
		c.Context().SetUserValue(logger.Endpoint, endpoint)

		// Propagate the context
		return c.Next()
	}
}

func CSRFProtection(cfg *config.Config) fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token", // Ensure this matches where you send the token
		CookieName:     "csrf_token",          // Name of the CSRF cookie
		CookieHTTPOnly: true,                  // Make the cookie HTTP-only
		CookieSecure:   cfg.Env == "production",
		CookieSameSite: "Strict",
		ContextKey:     "csrf_token", // Context key for the token
		Expiration:     cfg.JWTExpiresIn * time.Second,
	})
}
