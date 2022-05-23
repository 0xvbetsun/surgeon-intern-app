FROM golang:1.16.5-alpine3.13 AS builder
RUN apk --no-cache add git make build-base
ADD . /app
WORKDIR /app
RUN go mod download
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -a cmd/webservice/main.go

# final stage
FROM alpine:3.13
RUN apk --no-cache add ca-certificates yarn
RUN yarn global add wait-on
RUN mkdir configs
COPY --from=builder /app/main ./
COPY --from=builder /app/configs/casbin.conf ./configs/casbin.conf
ENV LISTEN_PORT 8080
EXPOSE 8080
CMD ["./main"]