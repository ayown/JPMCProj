# Multi-stage build for Go services
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build argument for service name
ARG SERVICE=api-gateway

# Build the specific service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/${SERVICE}/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/app .

# Copy migrations (if needed)
COPY --from=builder /build/internal/database/migrations ./migrations

# Expose port (will be overridden by docker-compose)
EXPOSE 8080

# Run the binary
CMD ["./app"]

