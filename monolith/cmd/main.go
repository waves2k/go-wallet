package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/waves2k/go-wallet/monolith/internal/config"
	"github.com/waves2k/go-wallet/monolith/internal/database"
	"github.com/waves2k/go-wallet/monolith/internal/logger"
	"github.com/waves2k/go-wallet/monolith/internal/user/handler"
	"github.com/waves2k/go-wallet/monolith/internal/user/repository"
	"github.com/waves2k/go-wallet/monolith/internal/user/service"
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

	userRepo := repository.NewPostgresqlUserRepository(pool)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	app := fiber.New()
	userHandler.InitRoutes(app)

	if err := app.Listen(config.ListenAddr); err != nil {
		log.Fatal(err)
	}
}
