.PHONY: up down test mocks

up:
	docker compose up -d

down:
	docker compose down

test:
	go test ./...

mocks:
	mockgen -source=internal/handlers/message_handler.go -destination=internal/mocks/mock_scan_processor.go -package=mocks
	mockgen -source=internal/consumer/scan_consumer.go -destination=internal/mocks/mock_message_handler.go -package=mocks
