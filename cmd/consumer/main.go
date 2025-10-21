package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/censys/scan-takehome/internal/repositories"
	"github.com/censys/scan-takehome/internal/workers"
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	subscriptionId := flag.String("subscription", "scan-sub", "GCP PubSub Subscription ID")
	dbHost := flag.String("db-host", getEnv("DB_HOST", "localhost"), "Database host")
	dbPort := flag.String("db-port", getEnv("DB_PORT", "5432"), "Database port")
	dbName := flag.String("db-name", getEnv("DB_NAME", "scans"), "Database name")
	dbUser := flag.String("db-user", getEnv("DB_USER", "postgres"), "Database user")
	dbPassword := flag.String("db-password", getEnv("DB_PASSWORD", "postgres"), "Database password")
	flag.Parse()

	if err := run(*projectId, *subscriptionId, *dbHost, *dbPort, *dbName, *dbUser, *dbPassword); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func run(projectId, subscriptionId, dbHost, dbPort, dbName, dbUser, dbPassword string) error {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	repo := repositories.NewPostgresRepository(db)

	config := workers.Config{
		ProjectID:      projectId,
		SubscriptionID: subscriptionId,
		Repository:     repo,
	}

	scanWorker, err := workers.NewScanWorker(config)
	if err != nil {
		return fmt.Errorf("failed to create scan worker: %w", err)
	}
	defer func() {
		if err := scanWorker.Stop(); err != nil {
			log.Printf("Error stopping scan worker: %v", err)
		}
	}()

	if err := scanWorker.Run(); err != nil {
		return fmt.Errorf("worker error: %w", err)
	}

	return nil
}
