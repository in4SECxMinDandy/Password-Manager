package vault

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/passwordmanager/backend/internal/common/errors"
	"github.com/passwordmanager/backend/internal/crypto"
)

type VaultService struct {
	vaultRepo  VaultRepository
	entryRepo  EntryRepository
	crypto     *crypto.CryptoService
}

func NewVaultService(
	vaultRepo VaultRepository,
	entryRepo EntryRepository,
	cryptoSvc *crypto.CryptoService,
) *VaultService {
	return &VaultService{
		vaultRepo:  vaultRepo,
		entryRepo:  entryRepo,
		crypto:     cryptoSvc,
	}
}

func (s *VaultService) CreateVault(ctx context.Context, userID uuid.UUID, req *CreateVaultRequest) (*Vault, error) {
	if req.Name == "" {
		return nil, errors.WithMessage(errors.ErrBadRequest, "vault name is required")
	}

	vault := &Vault{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.vaultRepo.Create(ctx, vault); err != nil {
		return nil, errors.ErrInternal
	}

	return vault, nil
}

func (s *VaultService) GetVault(ctx context.Context, userID, vaultID uuid.UUID) (*Vault, error) {
	vault, err := s.vaultRepo.GetByID(ctx, vaultID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	if vault.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return vault, nil
}

func (s *VaultService) ListVaults(ctx context.Context, userID uuid.UUID) (*VaultListResponse, error) {
	vaults, err := s.vaultRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	if vaults == nil {
		vaults = []Vault{}
	}

	return &VaultListResponse{
		Vaults: vaults,
		Total:  len(vaults),
	}, nil
}

func (s *VaultService) UpdateVault(ctx context.Context, userID, vaultID uuid.UUID, req *UpdateVaultRequest) (*Vault, error) {
	vault, err := s.GetVault(ctx, userID, vaultID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		vault.Name = req.Name
	}
	vault.UpdatedAt = time.Now()

	if err := s.vaultRepo.Update(ctx, vault); err != nil {
		return nil, errors.ErrInternal
	}

	return vault, nil
}

func (s *VaultService) DeleteVault(ctx context.Context, userID, vaultID uuid.UUID) error {
	_, err := s.GetVault(ctx, userID, vaultID)
	if err != nil {
		return err
	}

	if err := s.vaultRepo.Delete(ctx, vaultID); err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *VaultService) CreateEntry(ctx context.Context, userID, vaultID uuid.UUID, req *CreateEntryRequest) (*VaultEntry, error) {
	_, err := s.GetVault(ctx, userID, vaultID)
	if err != nil {
		return nil, err
	}

	if req.Title == "" {
		return nil, errors.WithMessage(errors.ErrBadRequest, "entry title is required")
	}

	entry := &VaultEntry{
		ID:        uuid.New(),
		VaultID:   vaultID,
		Type:      req.Type,
		Title:     req.Title,
		Data:      []byte(req.Data),
		Favorite:  req.Favorite,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.entryRepo.Create(ctx, entry); err != nil {
		return nil, errors.ErrInternal
	}

	return entry, nil
}

func (s *VaultService) GetEntry(ctx context.Context, userID, entryID uuid.UUID) (*VaultEntry, error) {
	entry, err := s.entryRepo.GetByID(ctx, entryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternal
	}

	vault, err := s.vaultRepo.GetByID(ctx, entry.VaultID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	if vault.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return entry, nil
}

func (s *VaultService) ListEntries(ctx context.Context, userID, vaultID uuid.UUID) (*EntryListResponse, error) {
	_, err := s.GetVault(ctx, userID, vaultID)
	if err != nil {
		return nil, err
	}

	entries, err := s.entryRepo.GetByVaultID(ctx, vaultID, false)
	if err != nil {
		return nil, errors.ErrInternal
	}

	if entries == nil {
		entries = []VaultEntry{}
	}

	return &EntryListResponse{
		Entries: entries,
		Total:   len(entries),
	}, nil
}

func (s *VaultService) UpdateEntry(ctx context.Context, userID, entryID uuid.UUID, req *UpdateEntryRequest) (*VaultEntry, error) {
	entry, err := s.GetEntry(ctx, userID, entryID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		entry.Title = req.Title
	}

	if req.Data != "" {
		entry.Data = []byte(req.Data)
	}

	if req.Favorite != nil {
		entry.Favorite = *req.Favorite
	}

	entry.UpdatedAt = time.Now()

	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, errors.ErrInternal
	}

	return entry, nil
}

func (s *VaultService) DeleteEntry(ctx context.Context, userID, entryID uuid.UUID) error {
	_, err := s.GetEntry(ctx, userID, entryID)
	if err != nil {
		return err
	}

	if err := s.entryRepo.Delete(ctx, entryID); err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (s *VaultService) ToggleFavorite(ctx context.Context, userID, entryID uuid.UUID) (*VaultEntry, error) {
	entry, err := s.GetEntry(ctx, userID, entryID)
	if err != nil {
		return nil, err
	}

	entry.Favorite = !entry.Favorite
	entry.UpdatedAt = time.Now()

	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, errors.ErrInternal
	}

	return entry, nil
}

func (s *VaultService) EncryptEntryData(data interface{}, key []byte) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return s.crypto.Encrypt(jsonData, key)
}

func (s *VaultService) DecryptEntryData(ciphertext []byte, key []byte) (json.RawMessage, error) {
	plaintext, err := s.crypto.Decrypt(ciphertext, key)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(plaintext), nil
}
