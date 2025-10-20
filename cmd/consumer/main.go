package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"github.com/censys/scan-takehome/pkg/scanning"
)

func main() {
	projectId := flag.String("project", "test-project", "GCP Project ID")
	subscriptionId := flag.String("subscription", "scan-sub", "GCP PubSub Subscription ID")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping consumer...")
		cancel()
	}()

	client, err := pubsub.NewClient(ctx, *projectId)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
	defer client.Close()

	subscription := client.Subscription(*subscriptionId)

	log.Printf("Starting consumer for subscription: %s", *subscriptionId)
	log.Println("Consumer is running... (Press Ctrl+C to stop)")

	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Received message ID: %s", msg.ID)

		var scan scanning.Scan
		if err := json.Unmarshal(msg.Data, &scan); err != nil {
			log.Printf("Failed to parse message: %v", err)
			msg.Nack()
			return
		}

		var serviceResponse string
		switch scan.DataVersion {
		case scanning.V1:
			if v1Data, ok := scan.Data.(*scanning.V1Data); ok {
				serviceResponse = string(v1Data.ResponseBytesUtf8)
			}
		case scanning.V2:
			if v2Data, ok := scan.Data.(*scanning.V2Data); ok {
				serviceResponse = v2Data.ResponseStr
			}
		}

		log.Printf("Processed scan: IP=%s, Port=%d, Service=%s, Response=%s",
			scan.Ip, scan.Port, scan.Service, serviceResponse)

		msg.Ack()
	})

	if err != nil {
		log.Printf("Subscription receive error: %v", err)
	}

	log.Println("Consumer stopped")
}
