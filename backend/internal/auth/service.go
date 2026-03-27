package auth

import (
	"context"
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/passwordmanager/backend/internal/common/cache"
	apperrors "github.com/passwordmanager/backend/internal/common/errors"
	"github.com/passwordmanager/backend/internal/crypto"
)

const (
	minPasswordLength    = 12
	maxFailedAttempts    = 5
	lockoutDuration      = 15 * time.Minute
)

type AuthService struct {
	userRepo       UserRepository
	cryptoService  *crypto.CryptoService
	tokenService   *TokenService
	cache          *cache.RedisClient
}

func NewAuthService(
	userRepo UserRepository,
	cryptoService *crypto.CryptoService,
	tokenService *TokenService,
	cache *cache.RedisClient,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		cryptoService: cryptoService,
		tokenService:  tokenService,
		cache:         cache,
	}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	if err := s.validateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := s.validatePassword(req.Password); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, apperrors.ErrEmailExists
	}

	salt, err := s.cryptoService.GenerateSalt()
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	passwordHash, err := s.cryptoService.DeriveKey([]byte(req.Password), salt)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	user := &User{
		ID:                uuid.New(),
		Email:             req.Email,
		MasterPasswordHash: passwordHash,
		PasswordSalt:      salt,
		MFAEnabled:        false,
		FailedLoginCount:  0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.ErrInternal
	}

	accessToken, refreshToken, err := s.tokenService.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:   900,
		User: &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			MFAEnabled: user.MFAEnabled,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	attempts, err := s.cache.GetLoginAttempts(ctx, req.Email)
	if err == nil && attempts >= maxFailedAttempts {
		return nil, apperrors.ErrAccountLocked
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.incrementFailedLogin(ctx, req.Email)
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, apperrors.ErrInternal
	}

	if user.IsLocked() {
		return nil, apperrors.ErrAccountLocked
	}

	passwordHash, err := s.cryptoService.DeriveKey([]byte(req.Password), user.PasswordSalt)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	if string(passwordHash) != string(user.MasterPasswordHash) {
		s.incrementFailedLogin(ctx, req.Email)
		user.IncrementFailedLogin()
		if user.FailedLoginCount >= maxFailedAttempts {
			user.Lock(lockoutDuration)
		}
		s.userRepo.Update(ctx, user)
		return nil, apperrors.ErrInvalidCredentials
	}

	s.cache.ResetLoginAttempts(ctx, req.Email)
	user.FailedLoginCount = 0
	user.UpdatedAt = time.Now()
	s.userRepo.Update(ctx, user)

	accessToken, refreshToken, err := s.tokenService.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:   900,
		User: &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			MFAEnabled: user.MFAEnabled,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.tokenService.RevokeAllUserTokens(userID)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	accessToken, newRefreshToken, err := s.tokenService.RefreshTokens(refreshToken)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:   900,
	}, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) incrementFailedLogin(ctx context.Context, email string) {
	s.cache.IncrementLoginAttempts(ctx, email)
}

func (s *AuthService) validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return apperrors.WithMessage(apperrors.ErrBadRequest, "invalid email format")
	}
	return nil
}

func (s *AuthService) validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return apperrors.WithMessage(apperrors.ErrPasswordTooWeak, "password must be at least 12 characters")
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return apperrors.WithMessage(apperrors.ErrPasswordTooWeak, "password must contain uppercase, lowercase, and digit")
	}

	return nil
}

func (s *AuthService) ValidatePasswordStrength(password string) (bool, string) {
	err := s.validatePassword(password)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			return false, appErr.Message
		}
		return false, "invalid password"
	}
	return true, ""
}
