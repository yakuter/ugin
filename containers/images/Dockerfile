# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /bin/ugin ./cmd/ugin/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /bin/ugin /app/ugin

# Copy config file
COPY config.yml /app/config.yml

# Expose port
EXPOSE 8081

# Run the application
ENTRYPOINT ["/app/ugin"]
