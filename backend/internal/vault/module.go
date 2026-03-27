package vault

import (
	"github.com/gofiber/fiber/v2"
	"github.com/passwordmanager/backend/internal/common/database"
	"github.com/passwordmanager/backend/internal/crypto"
)

type VaultModule struct {
	vaultService *VaultService
	handler     *VaultHandler
}

func NewVaultModule(db *database.PostgresDB, cryptoSvc *crypto.CryptoService) *VaultModule {
	vaultRepo := NewVaultRepository(db)
	entryRepo := NewEntryRepository(db)
	vaultService := NewVaultService(vaultRepo, entryRepo, cryptoSvc)
	handler := NewVaultHandler(vaultService)

	return &VaultModule{
		vaultService: vaultService,
		handler:     handler,
	}
}

func (m *VaultModule) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	vaults := api.Group("/vaults", authMiddleware)
	vaults.Get("/", m.handler.ListVaults)
	vaults.Post("/", m.handler.CreateVault)
	vaults.Get("/:id", m.handler.GetVault)
	vaults.Put("/:id", m.handler.UpdateVault)
	vaults.Delete("/:id", m.handler.DeleteVault)

	vaults.Get("/:id/entries", m.handler.ListEntries)
	vaults.Post("/:id/entries", m.handler.CreateEntry)

	entries := api.Group("/entries", authMiddleware)
	entries.Get("/:id", m.handler.GetEntry)
	entries.Put("/:id", m.handler.UpdateEntry)
	entries.Delete("/:id", m.handler.DeleteEntry)
	entries.Post("/:id/favorite", m.handler.ToggleFavorite)
}

func (m *VaultModule) GetVaultService() *VaultService {
	return m.vaultService
}
