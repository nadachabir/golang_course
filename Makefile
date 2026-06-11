# TaskFlow — developer task runner. Run `make help` to list targets.
# Recipe lines are TAB-indented (Make requires real tabs, not spaces).

-include .env
export

APP_NAME := taskflow
BIN_DIR  := bin
PKG      := ./...
COMPOSE  := docker compose -f docker/docker-compose.yml --env-file .env

.DEFAULT_GOAL := help

## help: list available targets
help:
	@grep -hE '^## ' $(MAKEFILE_LIST) | sed 's/## //'

## run: run the API server
run:
	go run ./cmd/api

## build: compile the API binary into ./bin
build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

## test: run all tests with the race detector
test:
	go test -race -count=1 $(PKG)

## cover: run tests and open a coverage summary
cover:
	go test -coverprofile=coverage.out $(PKG) && go tool cover -func=coverage.out

## tidy: sync go.mod and go.sum
tidy:
	go mod tidy

## fmt: format all Go code
fmt:
	gofmt -w .

## vet: report suspicious constructs
vet:
	go vet $(PKG)

## docker-up: start postgres + keycloak in the background
docker-up:
	$(COMPOSE) up -d

## docker-down: stop the local stack
docker-down:
	$(COMPOSE) down

## migrate-up: apply all DB migrations (needs golang-migrate CLI)
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

## migrate-down: roll back the most recent migration
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

.PHONY: help run build test cover tidy fmt vet docker-up docker-down migrate-up migrate-down
