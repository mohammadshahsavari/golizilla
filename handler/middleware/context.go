package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func SetTransaction(committer Committer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cm := committer.Begin()
		SetTx(c.UserContext(), cm)

		logger := GetLogger(c.UserContext())

		logger.Info("starting transaction")
		if err := c.Next(); err != nil {
			logger.Info("rollback on error", "error", err.Error())
			cm.Rollback()
			return err
		}

		err, ok := c.Locals(IsTxError).(error)
		if ok && err != nil {
			logger.Info("rollback on not ok response", "error", err.Error())
			cm.Rollback()
			return nil
		}

		if err := cm.Commit(); err != nil {
			logger.Info("commit error", "err", err.Error())
			cm.Rollback()
			return err
		}

		logger.Info("ending transaction")
		return nil
	}
}
