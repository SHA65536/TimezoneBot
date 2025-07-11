FROM golang:1.24.1-alpine3.21 AS builder

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate

RUN go build -o /timezonebot ./cmd/timezonebot
RUN go build -o /migrations ./cmd/migrations

FROM alpine:3.21 AS run
WORKDIR /app
RUN apk add --no-cache tzdata

USER 1001

COPY --from=builder /timezonebot ./timezonebot
COPY --from=builder /migrations ./migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./timezonebot"]