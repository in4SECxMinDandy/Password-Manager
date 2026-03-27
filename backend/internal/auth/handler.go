package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/passwordmanager/backend/internal/common/errors"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" {
		return errors.ErrBadRequest
	}

	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" {
		return errors.ErrBadRequest
	}

	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	if err := h.authService.Logout(c.Context(), userID); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "logged out successfully"})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.RefreshToken == "" {
		return errors.ErrBadRequest
	}

	response, err := h.authService.RefreshTokens(c.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	user, err := h.authService.GetUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		MFAEnabled: user.MFAEnabled,
		CreatedAt:  user.CreatedAt,
	})
}

func (h *AuthHandler) CheckPasswordStrength(c *fiber.Ctx) error {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	valid, message := h.authService.ValidatePasswordStrength(req.Password)
	return c.JSON(fiber.Map{
		"valid":   valid,
		"message": message,
	})
}
