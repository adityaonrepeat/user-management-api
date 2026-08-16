package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/adityaonrepeat/user-management-api/internal/repository"
	"github.com/adityaonrepeat/user-management-api/internal/service"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func ErrorHandler(log *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		message := "an unexpected error occurred"

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
			message = fiberErr.Message
		} else {
			log.Error("unhandled error",
				zap.Error(err),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
		}

		return c.Status(status).JSON(errorEnvelope{
			Error: errorBody{Code: codeFor(status), Message: message},
		})
	}
}

func codeFor(status int) string {
	switch status {
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "INTERNAL_ERROR"
	}
}

func mapServiceError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrInvalidDOB), errors.Is(err, service.ErrDOBInFuture):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return err
	}
}
