package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/passwordmanager/backend/internal/common/cache"
	"github.com/passwordmanager/backend/internal/common/config"
	"github.com/passwordmanager/backend/internal/common/database"
	"github.com/passwordmanager/backend/internal/crypto"
)

type AuthModule struct {
	userRepo      UserRepository
	authService   *AuthService
	tokenService  *TokenService
	jwtSecret     []byte
}

func NewAuthModule(
	db *database.PostgresDB,
	redis *cache.RedisClient,
	cryptoSvc *crypto.CryptoService,
	cfg config.JWTConfig,
) *AuthModule {
	userRepo := NewUserRepository(db)
	tokenService := NewTokenService(cfg, redis)
	authService := NewAuthService(userRepo, cryptoSvc, tokenService, redis)

	return &AuthModule{
		userRepo:     userRepo,
		authService:  authService,
		tokenService: tokenService,
		jwtSecret:    []byte(cfg.Secret),
	}
}

func (m *AuthModule) RegisterRoutes(api fiber.Router) {
	authHandler := NewAuthHandler(m.authService)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", m.GetAuthMiddleware(), authHandler.Logout)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Get("/me", m.GetAuthMiddleware(), authHandler.GetMe)
	auth.Post("/check-password", authHandler.CheckPasswordStrength)
}

func (m *AuthModule) GetAuthMiddleware() fiber.Handler {
	return jwtMiddleware.New(jwtMiddleware.Config{
		SigningKey: m.jwtSecret,
		SuccessHandler: func(c *fiber.Ctx) error {
			rawToken := c.Locals("user")
			if rawToken == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    "UNAUTHORIZED",
						"message": "Authentication required",
					},
				})
			}

			userID, email := extractJWTClaims(rawToken, m.jwtSecret)
			if userID == uuid.Nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": fiber.Map{
						"code":    "INVALID_TOKEN",
						"message": "Invalid token claims",
					},
				})
			}

			c.Locals("userID", userID)
			c.Locals("email", email)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
		},
	})
}

func extractJWTClaims(rawToken interface{}, secret []byte) (uuid.UUID, string) {
	val := fmt.Sprintf("%v", rawToken)

	parts := strings.Fields(val)
	if len(parts) < 2 {
		return uuid.Nil, ""
	}

	rawJWT := parts[0]
	if strings.HasPrefix(rawJWT, "&{") {
		rawJWT = rawJWT[2:]
	}
	if !strings.Contains(rawJWT, ".") {
		return uuid.Nil, ""
	}

	token, _, err := new(jwt.Parser).ParseUnverified(rawJWT, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ""
	}

	userIDStr, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)

	userID, _ := uuid.Parse(userIDStr)
	return userID, email
}

func (m *AuthModule) GetAuthService() *AuthService {
	return m.authService
}

func GetTokenFromHeader(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}
