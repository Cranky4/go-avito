FROM golang:1.18-alpine 

ENV GOOSE_DRIVER=
ENV GOOSE_DBSTRING=
ENV CODE_DIR=

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD goose -dir ${CODE_DIR} up
