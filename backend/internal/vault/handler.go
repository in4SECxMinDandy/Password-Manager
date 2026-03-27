package vault

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/passwordmanager/backend/internal/common/errors"
)

type VaultHandler struct {
	vaultService *VaultService
}

func NewVaultHandler(vaultService *VaultService) *VaultHandler {
	return &VaultHandler{vaultService: vaultService}
}

func (h *VaultHandler) CreateVault(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req CreateVaultRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	vault, err := h.vaultService.CreateVault(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(vault)
}

func (h *VaultHandler) GetVault(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	vaultID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	vault, err := h.vaultService.GetVault(c.Context(), userID, vaultID)
	if err != nil {
		return err
	}

	return c.JSON(vault)
}

func (h *VaultHandler) ListVaults(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	response, err := h.vaultService.ListVaults(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

func (h *VaultHandler) UpdateVault(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	vaultID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	var req UpdateVaultRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	vault, err := h.vaultService.UpdateVault(c.Context(), userID, vaultID, &req)
	if err != nil {
		return err
	}

	return c.JSON(vault)
}

func (h *VaultHandler) DeleteVault(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	vaultID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := h.vaultService.DeleteVault(c.Context(), userID, vaultID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *VaultHandler) CreateEntry(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	vaultID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	var req CreateEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	entry, err := h.vaultService.CreateEntry(c.Context(), userID, vaultID, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(entry)
}

func (h *VaultHandler) GetEntry(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	entry, err := h.vaultService.GetEntry(c.Context(), userID, entryID)
	if err != nil {
		return err
	}

	return c.JSON(entry)
}

func (h *VaultHandler) ListEntries(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	vaultID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	response, err := h.vaultService.ListEntries(c.Context(), userID, vaultID)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

func (h *VaultHandler) UpdateEntry(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	var req UpdateEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	entry, err := h.vaultService.UpdateEntry(c.Context(), userID, entryID, &req)
	if err != nil {
		return err
	}

	return c.JSON(entry)
}

func (h *VaultHandler) DeleteEntry(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := h.vaultService.DeleteEntry(c.Context(), userID, entryID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *VaultHandler) ToggleFavorite(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest
	}

	entry, err := h.vaultService.ToggleFavorite(c.Context(), userID, entryID)
	if err != nil {
		return err
	}

	return c.JSON(entry)
}
