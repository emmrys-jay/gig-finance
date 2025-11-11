package main

import (
	"flag"
	"log"
	"os"

	"github.com/emmrys-jay/gigmile/config"
	"github.com/emmrys-jay/gigmile/internal/migrations"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up, down, status")
	)
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	switch *command {
	case "up":
		log.Println("Running migrations up...")
		if err := migrations.RunMigrations(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations completed successfully")

	case "down":
		log.Println("Running migrations down...")
		if err := migrations.RunMigrationsDown(cfg); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations down completed successfully")

	case "status":
		log.Println("Checking migration status...")
		if err := migrations.GetMigrationStatus(cfg); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	default:
		log.Printf("Unknown command: %s\n", *command)
		log.Println("Usage: migrate -command=[up|down|status]")
		os.Exit(1)
	}
}
