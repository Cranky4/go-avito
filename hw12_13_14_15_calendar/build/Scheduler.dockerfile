# Собираем в гошке
FROM golang:1.18-alpine as builder

ENV BIN_FILE /opt/scheduler/scheduler-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.16.2

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/scheduler/scheduler-app"
COPY --from=builder ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/scheduler/config.toml
COPY ./configs/scheduler.toml ${CONFIG_FILE}

RUN mkdir logs

CMD ${BIN_FILE} -config ${CONFIG_FILE}
