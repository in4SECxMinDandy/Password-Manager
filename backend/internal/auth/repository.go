package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/passwordmanager/backend/internal/common/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type postgresUserRepository struct {
	db *database.PostgresDB
}

func NewUserRepository(db *database.PostgresDB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, master_password_hash, password_salt, mfa_enabled, mfa_secret, mfa_backup_codes, failed_login_count, locked_until, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.MasterPasswordHash,
		user.PasswordSalt,
		user.MFAEnabled,
		user.MFASecret,
		user.MFABackupCodes,
		user.FailedLoginCount,
		user.LockedUntil,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, master_password_hash, password_salt, mfa_enabled, mfa_secret, mfa_backup_codes, failed_login_count, locked_until, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &User{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.MasterPasswordHash,
		&user.PasswordSalt,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.MFABackupCodes,
		&user.FailedLoginCount,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, master_password_hash, password_salt, mfa_enabled, mfa_secret, mfa_backup_codes, failed_login_count, locked_until, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.MasterPasswordHash,
		&user.PasswordSalt,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.MFABackupCodes,
		&user.FailedLoginCount,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users SET 
			email = $2,
			master_password_hash = $3,
			password_salt = $4,
			mfa_enabled = $5,
			mfa_secret = $6,
			mfa_backup_codes = $7,
			failed_login_count = $8,
			locked_until = $9,
			updated_at = $10
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.MasterPasswordHash,
		user.PasswordSalt,
		user.MFAEnabled,
		user.MFASecret,
		user.MFABackupCodes,
		user.FailedLoginCount,
		user.LockedUntil,
		user.UpdatedAt,
	)
	return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
