package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/passwordmanager/backend/internal/common/cache"
	"github.com/passwordmanager/backend/internal/common/config"
)

type TokenService struct {
	jwtSecret      []byte
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
	redis          *cache.RedisClient
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func NewTokenService(cfg config.JWTConfig, redis *cache.RedisClient) *TokenService {
	return &TokenService{
		jwtSecret:     []byte(cfg.Secret),
		accessExpiry:  time.Duration(cfg.AccessExpiry) * time.Minute,
		refreshExpiry: time.Duration(cfg.RefreshExpiry) * 24 * time.Hour,
		redis:         redis,
	}
}

func (s *TokenService) GenerateTokens(userID uuid.UUID, email string) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	tokenID := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.redis.SetRefreshToken(ctx, tokenID, userID.String(), s.refreshExpiry); err != nil {
		return "", "", err
	}

	return accessToken, tokenID + ":" + refreshToken, nil
}

func (s *TokenService) generateAccessToken(userID uuid.UUID, email string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
		UserID: userID.String(),
		Email:  email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID.String(),
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *TokenService) RefreshTokens(refreshToken string) (newAccess, newRefresh string, err error) {
	parts := parseRefreshTokenParts(refreshToken)
	if len(parts) != 2 {
		return "", "", errors.New("invalid refresh token format")
	}

	tokenID, actualToken := parts[0], parts[1]

	ctx := context.Background()
	_, err = s.redis.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	token, err := jwt.ParseWithClaims(actualToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", "", errors.New("invalid user ID in token")
	}

	if err := s.redis.RevokeRefreshToken(ctx, tokenID); err != nil {
		return "", "", err
	}

	newAccess, newRefresh, err = s.GenerateTokens(userID, "")
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *TokenService) RevokeToken(tokenID string) error {
	ctx := context.Background()
	return s.redis.RevokeRefreshToken(ctx, tokenID)
}

func (s *TokenService) RevokeAllUserTokens(userID uuid.UUID) error {
	ctx := context.Background()
	return s.redis.RevokeAllUserTokens(ctx, userID.String())
}

func parseRefreshTokenParts(token string) []string {
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		return []string{}
	}
	return parts
}
