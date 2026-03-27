package errors

import (
	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrBadRequest = &AppError{Code: "BAD_REQUEST", Message: "Invalid request", StatusCode: fiber.StatusBadRequest}
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized", StatusCode: fiber.StatusUnauthorized}
	ErrForbidden = &AppError{Code: "FORBIDDEN", Message: "Access denied", StatusCode: fiber.StatusForbidden}
	ErrNotFound = &AppError{Code: "NOT_FOUND", Message: "Resource not found", StatusCode: fiber.StatusNotFound}
	ErrConflict = &AppError{Code: "CONFLICT", Message: "Resource already exists", StatusCode: fiber.StatusConflict}
	ErrInternal = &AppError{Code: "INTERNAL_ERROR", Message: "An internal error occurred", StatusCode: fiber.StatusInternalServerError}
	ErrTooManyRequests = &AppError{Code: "TOO_MANY_REQUESTS", Message: "Too many requests, please try again later", StatusCode: fiber.StatusTooManyRequests}
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password", StatusCode: fiber.StatusUnauthorized}
	ErrAccountLocked = &AppError{Code: "ACCOUNT_LOCKED", Message: "Account is temporarily locked", StatusCode: fiber.StatusForbidden}
	ErrInvalidToken = &AppError{Code: "INVALID_TOKEN", Message: "Invalid or expired token", StatusCode: fiber.StatusUnauthorized}
	ErrPasswordTooWeak = &AppError{Code: "PASSWORD_TOO_WEAK", Message: "Password does not meet security requirements", StatusCode: fiber.StatusBadRequest}
	ErrEmailExists = &AppError{Code: "EMAIL_EXISTS", Message: "Email already registered", StatusCode: fiber.StatusConflict}
)

func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: statusCode}
}

func WithMessage(err *AppError, message string) *AppError {
	return &AppError{Code: err.Code, Message: message, StatusCode: err.StatusCode}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok {
		return c.Status(appErr.StatusCode).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "HTTP_ERROR",
				"message": fiberErr.Message,
			},
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		},
	})
}
