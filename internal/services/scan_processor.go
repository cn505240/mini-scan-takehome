package services

import (
	"context"
	"log"

	"github.com/censys/scan-takehome/internal/domain"
)

type ScanRepository interface {
	GetLatestScan(ctx context.Context, ip string, port uint32, service string) (*domain.ServiceScan, error)
	UpsertScan(ctx context.Context, scan *domain.ServiceScan) error
}

type ScanProcessor struct {
	repository ScanRepository
}

func NewScanProcessor(repository ScanRepository) *ScanProcessor {
	return &ScanProcessor{
		repository: repository,
	}
}

func (sp *ScanProcessor) ProcessScanResult(ctx context.Context, scan *domain.ServiceScan) error {
	latestScan, err := sp.repository.GetLatestScan(ctx, scan.IP, scan.Port, scan.Service)
	if err != nil {
		log.Printf("Failed to get latest scan for %s:%d/%s: %v",
			scan.IP, scan.Port, scan.Service, err)
		return err
	}

	if latestScan != nil && !scan.IsNewerThan(latestScan.LastScanned) {
		log.Printf("Ignoring older scan for %s:%d/%s (latest: %v, received: %v)",
			scan.IP, scan.Port, scan.Service,
			latestScan.LastScanned, scan.LastScanned)
		return nil
	}

	if err := sp.repository.UpsertScan(ctx, scan); err != nil {
		log.Printf("Failed to upsert scan for %s:%d/%s: %v",
			scan.IP, scan.Port, scan.Service, err)
		return err
	}

	log.Printf("Updated scan for %s:%d/%s with timestamp %v",
		scan.IP, scan.Port, scan.Service, scan.LastScanned)
	return nil
}
