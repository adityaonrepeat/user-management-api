package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/adityaonrepeat/user-management-api/config"
	"github.com/adityaonrepeat/user-management-api/internal/logger"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("invalid configuration", zap.Error(err))
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to create connection pool", zap.Error(err))
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		log.Fatal("database unreachable", zap.Error(err))
	}
	log.Info("connected to database")

	app := fiber.New(fiber.Config{
		AppName:               "user-management-api",
		DisableStartupMessage: true,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		if err := pool.Ping(c.UserContext()); err != nil {
			log.Error("health check failed", zap.Error(err))
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	go func() {
		addr := ":" + cfg.ServerPort
		log.Info("server listening", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down")
	if err := app.Shutdown(); err != nil {
		log.Error("shutdown failed", zap.Error(err))
	}
}
