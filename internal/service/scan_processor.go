package service

import (
	"context"
	"log"

	"github.com/censys/scan-takehome/internal/domain"
)

type ScanProcessor struct{}

func NewScanProcessor() *ScanProcessor {
	return &ScanProcessor{}
}

func (sp *ScanProcessor) ProcessScanResult(ctx context.Context, scanResult *domain.ScanResult) error {
	log.Printf("Processing scan: %s:%d/%s - %s",
		scanResult.IP, scanResult.Port, scanResult.Service, scanResult.Response)
	return nil
}
