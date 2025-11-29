# Builder stage
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/server

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app

RUN useradd -u 10001 -r -s /bin/false appuser && \
    apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/.env.example /app/.env.example
COPY --from=builder /app/docker-entrypoint.sh /app/docker-entrypoint.sh

ENV APP_ENV=production
EXPOSE 8080
USER appuser

ENTRYPOINT ["/bin/bash", "/app/docker-entrypoint.sh"]
