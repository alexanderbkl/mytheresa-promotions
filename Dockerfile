# syntax=docker/dockerfile:1

# Build Stage
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Install Git (if needed for dependencies)
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy the built binary and products.json
COPY --from=builder /app/main .
COPY --from=builder /app/products.json .

# Expose the application port
EXPOSE 8080

# Set the entry point
CMD ["./main"]
