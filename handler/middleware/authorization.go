package middleware

import (
	"golizilla/handler/presenter"
	"golizilla/internal/apperrors"
	"golizilla/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuthorizationMiddleware(authoriztionService service.IAuthorizationService) func(requiredPrivileges ...string) fiber.Handler {
	return func(requiredPrivileges ...string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			ctx := c.Context()

			hasPrivilege, err := authoriztionService.IsAuthorized(ctx, c.UserContext(), c.Locals("user_id").(uuid.UUID), requiredPrivileges...)

			if err != nil {
				// log
				return err
			}

			if !hasPrivilege {
				return presenter.SendError(c, fiber.StatusForbidden, apperrors.ErrLackOfAuthorization.Error())
			}

			return c.Next()
		}
	}
}
