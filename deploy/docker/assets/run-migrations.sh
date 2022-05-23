#!/bin/sh

set -e

goose -dir=. postgres "$DATABASE_URL" up
goose -dir=. postgres "$DATABASE_URL" status
