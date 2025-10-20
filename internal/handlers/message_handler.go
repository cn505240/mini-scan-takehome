package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/censys/scan-takehome/internal/domain"
	"github.com/censys/scan-takehome/pkg/scanning"
)

type ScanProcessor interface {
	ProcessScanResult(ctx context.Context, scanResult domain.ScanResult) error
}

type MessageHandler struct {
	processor ScanProcessor
}

func NewMessageHandler(processor ScanProcessor) *MessageHandler {
	return &MessageHandler{
		processor: processor,
	}
}

func (mh *MessageHandler) HandleMessage(ctx context.Context, msgData []byte) error {
	var rawScan scanning.Scan
	if err := json.Unmarshal(msgData, &rawScan); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return err
	}

	scanResult, err := domain.ConvertScanToDomain(rawScan)
	if err != nil {
		log.Printf("Failed to convert message to domain model: %v", err)
		return err
	}

	return mh.processor.ProcessScanResult(ctx, scanResult)
}
