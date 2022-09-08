FROM golang:1.18-alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/calendar ./cmd/calendar

# -----------------------------------------

FROM alpine:3.16.2 as calendar

WORKDIR /usr/src/app

COPY --from=builder /usr/local/bin/calendar /usr/src/app

RUN mkdir logs

CMD ["/usr/src/app/calendar", "-config", "./config.toml"]