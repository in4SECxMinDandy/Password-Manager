package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/passwordmanager/backend/internal/auth"
	"github.com/passwordmanager/backend/internal/common/config"
	"github.com/passwordmanager/backend/internal/common/database"
	"github.com/passwordmanager/backend/internal/common/cache"
	"github.com/passwordmanager/backend/internal/common/errors"
	"github.com/passwordmanager/backend/internal/common/middleware"
	"github.com/passwordmanager/backend/internal/crypto"
	"github.com/passwordmanager/backend/internal/vault"
	"github.com/rs/zerolog"
)

func main() {
	logger := setupLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()
	
	db, err := database.NewPostgresDB(ctx, cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	logger.Info().Msg("Connected to PostgreSQL")

	if err := db.RunMigrations(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}
	logger.Info().Msg("Database migrations completed")

	redisClient, err := cache.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()
	logger.Info().Msg("Connected to Redis")

	cryptoService := crypto.NewCryptoService()
	
	authModule := auth.NewAuthModule(db, redisClient, cryptoService, cfg.JWT)
	vaultModule := vault.NewVaultModule(db, cryptoService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error().
				Err(err).
				Str("path", c.Path()).
				Int("status", c.Response().StatusCode()).
				Msg("Request error")
			return errors.ErrorHandler(c, err)
		},
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Server.AllowOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(logger))

	api := app.Group("/api/v1")

	authModule.RegisterRoutes(api)
	vaultModule.RegisterRoutes(api, authModule.GetAuthMiddleware())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := cfg.Server.Address
		logger.Info().Str("address", addr).Msg("Starting server")
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	logger.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server stopped")
}

func setupLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	l := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Enable console writer with colors for Windows
	return l.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
