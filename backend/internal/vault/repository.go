package vault

import (
	"context"

	"github.com/google/uuid"
	"github.com/passwordmanager/backend/internal/common/database"
)

type VaultRepository interface {
	Create(ctx context.Context, vault *Vault) error
	GetByID(ctx context.Context, id uuid.UUID) (*Vault, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Vault, error)
	Update(ctx context.Context, vault *Vault) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type EntryRepository interface {
	Create(ctx context.Context, entry *VaultEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*VaultEntry, error)
	GetByVaultID(ctx context.Context, vaultID uuid.UUID, includeDeleted bool) ([]VaultEntry, error)
	Update(ctx context.Context, entry *VaultEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

type postgresVaultRepository struct {
	db *database.PostgresDB
}

func NewVaultRepository(db *database.PostgresDB) VaultRepository {
	return &postgresVaultRepository{db: db}
}

func (r *postgresVaultRepository) Create(ctx context.Context, vault *Vault) error {
	query := `
		INSERT INTO vaults (id, user_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		vault.ID,
		vault.UserID,
		vault.Name,
		vault.CreatedAt,
		vault.UpdatedAt,
	)
	return err
}

func (r *postgresVaultRepository) GetByID(ctx context.Context, id uuid.UUID) (*Vault, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM vaults WHERE id = $1
	`
	vault := &Vault{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&vault.ID,
		&vault.UserID,
		&vault.Name,
		&vault.CreatedAt,
		&vault.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return vault, nil
}

func (r *postgresVaultRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Vault, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM vaults WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vaults []Vault
	for rows.Next() {
		var vault Vault
		if err := rows.Scan(
			&vault.ID,
			&vault.UserID,
			&vault.Name,
			&vault.CreatedAt,
			&vault.UpdatedAt,
		); err != nil {
			return nil, err
		}
		vaults = append(vaults, vault)
	}

	return vaults, nil
}

func (r *postgresVaultRepository) Update(ctx context.Context, vault *Vault) error {
	query := `
		UPDATE vaults SET name = $2, updated_at = $3
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, vault.ID, vault.Name, vault.UpdatedAt)
	return err
}

func (r *postgresVaultRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vaults WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *postgresVaultRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM vaults WHERE user_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

type postgresEntryRepository struct {
	db *database.PostgresDB
}

func NewEntryRepository(db *database.PostgresDB) EntryRepository {
	return &postgresEntryRepository{db: db}
}

func (r *postgresEntryRepository) Create(ctx context.Context, entry *VaultEntry) error {
	query := `
		INSERT INTO vault_entries (id, vault_id, type, title, data, favorite, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		entry.ID,
		entry.VaultID,
		entry.Type,
		entry.Title,
		entry.Data,
		entry.Favorite,
		entry.CreatedAt,
		entry.UpdatedAt,
	)
	return err
}

func (r *postgresEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*VaultEntry, error) {
	query := `
		SELECT id, vault_id, type, title, data, favorite, deleted_at, created_at, updated_at
		FROM vault_entries WHERE id = $1
	`
	entry := &VaultEntry{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&entry.ID,
		&entry.VaultID,
		&entry.Type,
		&entry.Title,
		&entry.Data,
		&entry.Favorite,
		&entry.DeletedAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *postgresEntryRepository) GetByVaultID(ctx context.Context, vaultID uuid.UUID, includeDeleted bool) ([]VaultEntry, error) {
	query := `
		SELECT id, vault_id, type, title, data, favorite, deleted_at, created_at, updated_at
		FROM vault_entries WHERE vault_id = $1
	`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []VaultEntry
	for rows.Next() {
		var entry VaultEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.VaultID,
			&entry.Type,
			&entry.Title,
			&entry.Data,
			&entry.Favorite,
			&entry.DeletedAt,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *postgresEntryRepository) Update(ctx context.Context, entry *VaultEntry) error {
	query := `
		UPDATE vault_entries SET title = $2, data = $3, favorite = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, entry.ID, entry.Title, entry.Data, entry.Favorite, entry.UpdatedAt)
	return err
}

func (r *postgresEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE vault_entries SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *postgresEntryRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vault_entries WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *postgresEntryRepository) Restore(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE vault_entries SET deleted_at = NULL WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
