APP=api

.PHONY: run dev test lint build up down migrate seed

dev:
	go run ./cmd/api

build:
	go build -o bin/$(APP) ./cmd/api

test:
	go test ./... -race -cover

up:
	docker compose up -d --build

down:
	docker compose down -v

migrate:
	go run ./cmd/api --migrate

seed:
	go run ./cmd/api --seed
