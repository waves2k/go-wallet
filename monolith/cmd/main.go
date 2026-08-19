package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/waves2k/go-wallet/monolith/internal/config"
	"github.com/waves2k/go-wallet/monolith/internal/database"
	"github.com/waves2k/go-wallet/monolith/internal/logger"
	"github.com/waves2k/go-wallet/monolith/internal/user/handler"
	userRepo "github.com/waves2k/go-wallet/monolith/internal/user/repository"
	walRepo "github.com/waves2k/go-wallet/monolith/internal/wallet/repository"

	userSvc "github.com/waves2k/go-wallet/monolith/internal/user/service"
	walSvc "github.com/waves2k/go-wallet/monolith/internal/wallet/service"
)

func main() {
	config := config.LoadConfig()
	logger.InitLogger()

	ctx := context.Background()

	pool, err := database.ConnectWithRetry(ctx, config.GetConnectionString())
	if err != nil {
		logger.Log.Error("Criticaly Error: Could not connect to database", err)
	}
	defer pool.Close()

	walRepo := walRepo.NewPostgresqlUserRepository(pool)
	userRepo := userRepo.NewPostgresqlUserRepository(pool, walRepo)

	userSvc := userSvc.NewUserService(userRepo)
	walSvc := walSvc.NewWalletService(walRepo)

	userHandler := handler.NewUserHandler(userSvc, walSvc)

	app := fiber.New()
	userHandler.InitRoutes(app)

	if err := app.Listen(config.ListenAddr); err != nil {
		log.Fatal(err)
	}
}
