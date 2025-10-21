.PHONY: up down test mocks lint lint-fix migrate

up:
	docker compose up -d

down:
	docker compose down

test:
	go test ./...

migrate:
	PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d scans -f internal/db/migrations/001_create_service_records.sql

mocks:
	mockgen -source=internal/handlers/message_handler.go -destination=internal/mocks/mock_scan_processor.go -package=mocks
	mockgen -source=internal/workers/scan_worker.go -destination=internal/mocks/mock_message_handler.go -package=mocks
	mockgen -source=internal/services/scan_processor.go -destination=internal/mocks/mock_scan_repository.go -package=mocks

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
