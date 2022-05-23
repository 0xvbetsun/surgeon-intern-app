FROM golang:1.16.5-alpine3.13 AS builder
RUN apk add --no-cache git make build-base
WORKDIR /app
RUN go get -u github.com/pressly/goose/cmd/goose
RUN which goose
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -a cmd/seeder/main.go

# final stage
FROM alpine:3.13
RUN apk --no-cache add ca-certificates yarn
RUN yarn global add wait-on
COPY --from=builder /app/main ./
COPY --from=builder /go/bin/goose /usr/bin/goose
COPY --from=builder /app/deploy/goose-migrations /migrations
COPY --from=builder /app/surgeries.csv /surgeries.csv
COPY --from=builder /app/deploy/base-seed/base-seed.json /base-seed.json
COPY deploy/docker/assets/seeder.sh /seeder.sh
RUN chmod +x /seeder.sh
ENTRYPOINT ["/seeder.sh"]