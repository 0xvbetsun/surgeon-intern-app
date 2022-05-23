FROM golang:1.16.5-alpine3.13 AS builder
RUN apk add --no-cache git build-base
WORKDIR /app
COPY deploy/goose-migrations .
COPY deploy/docker/assets/run-migrations.sh .
RUN chmod +x run-migrations.sh
RUN go get -u github.com/pressly/goose/cmd/goose
