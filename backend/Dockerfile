# ===== Build stage =====
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN go build -o server .
RUN go build -o redis-consumer ./cmd/redis-consumer

# ===== Runtime stage =====
FROM alpine:latest

# Add CA certificates
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/server .
COPY --from=builder /app/redis-consumer .

# Copy start script
COPY start.sh .

# Make sure start.sh is executable
RUN chmod +x start.sh

# Expose port
EXPOSE 8080

# Start both server and consumer
CMD ["sh", "start.sh"]
