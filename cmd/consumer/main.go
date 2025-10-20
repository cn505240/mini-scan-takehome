package main

import (
	"flag"
	"log"

	"github.com/censys/scan-takehome/internal/consumer"
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	subscriptionId := flag.String("subscription", "scan-sub", "GCP PubSub Subscription ID")
	flag.Parse()

	// Create worker
	config := consumer.Config{
		ProjectID:      *projectId,
		SubscriptionID: *subscriptionId,
	}

	scanWorker, err := consumer.NewScanWorker(config)
	if err != nil {
		log.Fatalf("Failed to create scan worker: %v", err)
	}
	defer scanWorker.Stop()

	// Run the worker
	if err := scanWorker.Run(); err != nil {
		log.Printf("Worker error: %v", err)
	}
}
