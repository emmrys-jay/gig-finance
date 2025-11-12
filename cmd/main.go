package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emmrys-jay/gigmile/config"
	"github.com/emmrys-jay/gigmile/internal/cache"
	"github.com/emmrys-jay/gigmile/internal/database"
	"github.com/emmrys-jay/gigmile/internal/repository"
	"github.com/emmrys-jay/gigmile/internal/router"
	"github.com/emmrys-jay/gigmile/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Initialize repositories
	customerRepo := repository.NewCustomerRepository(db.Pool)
	accountRepo := repository.NewAccountRepository(db.Pool)
	transactionRepo := repository.NewTransactionRepository(db.Pool)

	// Initialize services
	customerService := service.NewCustomerService(customerRepo, accountRepo)
	paymentService := service.NewPaymentService(customerRepo, accountRepo, transactionRepo, redisCache)
	deploymentService := service.NewDeploymentService(customerRepo, accountRepo, transactionRepo, redisCache)
	transactionService := service.NewTransactionService(transactionRepo)
	accountService := service.NewAccountService(accountRepo)

	// Initialize router
	r := router.NewRouter(customerService, paymentService, deploymentService, transactionService, accountService)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("API endpoints: http://localhost%s/api/v1", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
