#!/bin/sh

set -e

wait-on tcp:db:5432
#goose -dir=deploy/goose-migrations postgres "user=${DB_USER} host=${DB_HOST} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable" reset
goose -dir=/migrations postgres "user=${DB_USER} host=${DB_HOST} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable" up
goose -dir=/migrations postgres "user=${DB_USER} host=${DB_HOST} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable" status
export DB_CONN_STRING="user=${DB_USER} host=${DB_HOST} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable"

./main
