BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

DATABASE_DOCKER_IMG_NAME="calendar-postgres"
DATABASE_NAME="calendar"
DATABASE_USER="calendar_user"
DATABASE_PASSWORD="secret"
DATABASE_PORT="55432"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
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

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.48

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

build-database-img:
	docker run \
		--name $(DATABASE_DOCKER_IMG_NAME) \
		-e POSTGRES_DB=${DATABASE_NAME} \
		-e POSTGRES_USER=${DATABASE_USER} \
		-e POSTGRES_PASSWORD=${DATABASE_PASSWORD} \
		-p "${DATABASE_PORT}:5432" \
		postgres:14-alpine
start-database-img:
	docker start ${DATABASE_DOCKER_IMG_NAME}
stop-database-img:
	docker stop ${DATABASE_DOCKER_IMG_NAME}

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