.PHONY: run build tidy docker-up docker-down

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down
