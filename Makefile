.PHONY: swagger build test run docker-up docker-down

swagger:
	swag init -g cmd/api/main.go -o openapi/docs

diagrams:
	plantuml docs/cjm/*.puml

lint:
	golangci-lint run ./...

build:
	go build ./...

test:
	go test ./...

run:
	go run ./cmd/api

docker-up:
	docker compose up --build

docker-down:
	docker compose down