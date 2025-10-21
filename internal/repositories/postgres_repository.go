package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/censys/scan-takehome/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetLatestScan(ctx context.Context, ip string, port uint32, service string) (*domain.ServiceScan, error) {
	query := `
		SELECT ip, port, service, response, last_scanned 
		FROM service_scans 
		WHERE ip = $1 AND port = $2 AND service = $3
		ORDER BY last_scanned DESC 
		LIMIT 1`

	var scan domain.ServiceScan
	err := r.db.QueryRowContext(ctx, query, ip, port, service).Scan(
		&scan.IP, &scan.Port, &scan.Service, &scan.Response, &scan.LastScanned,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest scan: %w", err)
	}

	return &scan, nil
}

func (r *PostgresRepository) UpsertScan(ctx context.Context, scan *domain.ServiceScan) error {
	query := `
		INSERT INTO service_scans (ip, port, service, response, last_scanned)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ip, port, service) 
		DO UPDATE SET 
			response = EXCLUDED.response,
			last_scanned = EXCLUDED.last_scanned
		WHERE service_scans.last_scanned < EXCLUDED.last_scanned`

	_, err := r.db.ExecContext(ctx, query,
		scan.IP, scan.Port, scan.Service, scan.Response, scan.LastScanned)

	if err != nil {
		return fmt.Errorf("failed to upsert scan: %w", err)
	}

	return nil
}
