# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Stage 2: Production
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Copy any additional configuration files if needed
COPY --from=builder /app/configs ./configs

# Create non-root user
RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["./main"]