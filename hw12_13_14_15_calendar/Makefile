BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

SCHEDULER_BIN := "./bin/scheduler"

SENDER_BIN := "./bin/sender"

DATABASE_DOCKER_IMG_NAME="calendar-postgres"
DATABASE_NAME="calendar"
DATABASE_USER="calendar_user"
DATABASE_PASSWORD="secret"
DATABASE_PORT="55432"

DATABASE_DOCKER_IMG_NAME="calendar-postgres"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: build-calendar build-scheduler build-sender

# Calendar
build-calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run-calendar: build-calendar
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

# Scheduler
build-scheduler:
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler

run-scheduler: build-scheduler
	$(SCHEDULER_BIN) -config ./configs/scheduler.toml

# Sender
build-sender:
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run-sender: build-sender
	$(SENDER_BIN) -config ./configs/sender.toml

# Quallity
test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.48.0

lint: install-lint-deps
	golangci-lint run ./...

# Database
install-goose:
	(which goose > /dev/null) || go install github.com/pressly/goose/v3/cmd/goose@latest
migrate-status: install-goose
	 goose -dir ./migrations postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} port=${DATABASE_PORT} sslmode=disable" status

migrate-create: install-goose
	 goose -dir ./migrations postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} port=${DATABASE_PORT} sslmode=disable" create calendar sql

migrate: install-goose
	 goose -dir ./migrations postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} port=${DATABASE_PORT} sslmode=disable" up

migrate-down: install-goose	
	 goose -dir ./migrations postgres "user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} port=${DATABASE_PORT} sslmode=disable" down

install-protobuf:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

# GRPC
generate:
	protoc ./api/EventService.proto --go_out=./api --go-grpc_out=./api

# Dev
up-dev:
	DATABASE_NAME="calendar" \
	DATABASE_USER="calendar_user" \
	DATABASE_PASSWORD="secret" \
	DATABASE_PORT="55432" \
	ZOOKEEPER_PORT="2181" \
	KAFKA_PORT="9093" \
	KAFKA_UI_PORT="22181" \
	docker-compose -f deployments/docker-compose.dev.yaml up -d

down-dev:
	docker-compose -f deployments/docker-compose.dev.yaml down
logs-dev:
	docker-compose -f deployments/docker-compose.dev.yaml logs calendar scheduler sender migrator -f

.PHONY: build run build-img run-img version test lint
