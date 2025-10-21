package workers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"

	"github.com/censys/scan-takehome/internal/handlers"
	"github.com/censys/scan-takehome/internal/services"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, msgData []byte) error
}

type Config struct {
	ProjectID      string
	SubscriptionID string
	Repository     services.ScanRepository
}

type ScanWorker struct {
	config         Config
	client         *pubsub.Client
	subscription   *pubsub.Subscription
	messageHandler MessageHandler
}

func NewScanWorker(config Config) (*ScanWorker, error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}

	subscription := client.Subscription(config.SubscriptionID)

	processor := services.NewScanProcessor(config.Repository)
	messageHandler := handlers.NewMessageHandler(processor)

	return &ScanWorker{
		config:         config,
		client:         client,
		subscription:   subscription,
		messageHandler: messageHandler,
	}, nil
}

func (sw *ScanWorker) Start(ctx context.Context) error {
	log.Printf("Starting worker for subscription: %s", sw.config.SubscriptionID)
	log.Println("Worker is running... (Press Ctrl+C to stop)")

	return sw.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Received message ID: %s", msg.ID)

		if err := sw.messageHandler.HandleMessage(ctx, msg.Data); err != nil {
			log.Printf("Failed to process message: %v", err)
			msg.Nack()
			return
		}

		msg.Ack()
	})
}

func (sw *ScanWorker) Stop() error {
	return sw.client.Close()
}

func (sw *ScanWorker) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping worker...")
		cancel()
	}()

	err := sw.Start(ctx)
	if err != nil {
		return err
	}

	log.Println("Worker stopped")
	return nil
}
