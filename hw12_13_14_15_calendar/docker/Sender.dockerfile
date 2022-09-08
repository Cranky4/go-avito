FROM golang:1.18-alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ENV GIT_HASH = $(shell git log --format="%h" -n 1)
ENV LDFLAGS = -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

COPY . .
RUN go build -v -o /usr/local/bin/sender -ldflags "$(LDFLAGS)" ./cmd/sender

# -----------------------------------------

FROM alpine:3.16.2 as sender

WORKDIR /usr/src/app

COPY --from=builder /usr/local/bin/sender /usr/src/app

RUN mkdir logs

CMD ["/usr/src/app/sender", "-config", "./config.toml"]