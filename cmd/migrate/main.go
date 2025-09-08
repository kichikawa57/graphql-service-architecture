package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"graphql-backend/database"
	"graphql-backend/migrations"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var (
		action     = flag.String("action", "up", "Migration action: up, down, status, create")
		steps      = flag.Int("steps", 1, "Number of steps for down migration")
		name       = flag.String("name", "", "Name for new migration (required for create action)")
		migrateDir = flag.String("dir", "migrations", "Directory containing migration files")
	)
	flag.Parse()

	// Get absolute path for migrations directory
	backendDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}
	migrationsDir := filepath.Join(backendDir, *migrateDir)

	switch *action {
	case "create":
		if *name == "" {
			log.Fatal("Migration name is required for create action. Use -name flag")
		}
		if err := migrations.CreateMigration(migrationsDir, *name); err != nil {
			log.Fatal("Failed to create migration:", err)
		}
		return

	case "up", "down", "status":
		// Connect to database for other actions
		db, err := database.NewDB()
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		migrator := migrations.NewMigrator(db.DB)

		switch *action {
		case "up":
			if err := migrator.Up(migrationsDir); err != nil {
				log.Fatal("Failed to run migrations:", err)
			}

		case "down":
			if err := migrator.Down(migrationsDir, *steps); err != nil {
				log.Fatal("Failed to rollback migrations:", err)
			}

		case "status":
			if err := migrator.Status(migrationsDir); err != nil {
				log.Fatal("Failed to get migration status:", err)
			}
		}

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status, create")
		flag.Usage()
		os.Exit(1)
	}
}
