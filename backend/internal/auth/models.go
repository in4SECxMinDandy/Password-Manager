package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	MasterPasswordHash []byte    `json:"-"`
	PasswordSalt      []byte     `json:"-"`
	MFAEnabled        bool       `json:"mfa_enabled"`
	MFASecret         []byte     `json:"-"`
	MFABackupCodes    [][]byte   `json:"-"`
	FailedLoginCount  int        `json:"failed_login_count"`
	LockedUntil       *time.Time `json:"locked_until,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

func (u *User) Lock(duration time.Duration) {
	lockedUntil := time.Now().Add(duration)
	u.LockedUntil = &lockedUntil
}

func (u *User) Unlock() {
	u.LockedUntil = nil
	u.FailedLoginCount = 0
}

func (u *User) IncrementFailedLogin() {
	u.FailedLoginCount++
}

type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type RegisterRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn   int    `json:"expires_in"`
	User        *UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	MFAEnabled bool     `json:"mfa_enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type MFASetupResponse struct {
	Secret     string `json:"secret"`
	QRCodeURL  string `json:"qr_code_url"`
}

type MFAVerifyRequest struct {
	Code string `json:"code"`
}

type MFAEnableRequest struct {
	Code        string `json:"code"`
	BackupCodes []string `json:"backup_codes,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
