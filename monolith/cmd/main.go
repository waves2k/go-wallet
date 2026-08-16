package main

import (
	"context"
	"log"

	"github.com/waves2k/go-wallet/monolith/internal/config"
	"github.com/waves2k/go-wallet/monolith/internal/database"
)

func main() {
	config := config.LoadConfig()

	ctx := context.Background()

	pool, err := database.ConnectWithRetry(ctx, config.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

}
