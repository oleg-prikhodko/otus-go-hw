BIN_DIR := "./bin"
BINARIES := calendar scheduler sender
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER := postgres-calendar

#Rabbit
RABBIT_USER ?= rabbit
RABBIT_PASSWORD ?= password
RABBIT_PORT ?= 5672
RABBIT_PORT_UI ?= 15672
RABBIT_CONTAINER := rabbitmq


run-postgres:
	docker run -d --name $(POSTGRES_CONTAINER) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_PORT):5432 \
	-v postgres-data:/var/lib/postgresql/data \
	postgres:latest

run-rabbit:
	docker run -d --name $(RABBIT_CONTAINER) \
	-p $(RABBIT_PORT):5672 \
	-p $(RABBIT_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
	rabbitmq:management

build:
	mkdir -p $(BIN_DIR)
	$(foreach binary,$(BINARIES),go build -v -o $(BIN_DIR)/$(binary) -ldflags "$(LDFLAGS)" ./cmd/$(binary) &&) true

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...
test-integration:
	go test -v -tags=integration ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./cmd/... ./internal/...

migrate-up:
	goose up

migrate-down:
	goose down

generate:
	go generate ./internal/server/grpc/...

up:
	docker compose -f ./deployments/docker-compose.yaml --env-file ./deployments/.env.example up -d

down:
	docker compose -f ./deployments/docker-compose.yaml --env-file ./deployments/.env.example down

integration-tests:
	@docker compose -f ./deployments/docker-compose.tests.yaml --env-file ./deployments/.env.test.example up -d calendar scheduler sender; \
    	docker compose -f ./deployments/docker-compose.tests.yaml --env-file ./deployments/.env.test.example run --rm tests; \
    	status=$$?; \
    	docker compose -f ./deployments/docker-compose.tests.yaml --env-file ./deployments/.env.test.example down; \
    	exit $$status

.PHONY: build run build-img run-img version test test-integration lint migrate-up migrate-down generate up down integration-tests
