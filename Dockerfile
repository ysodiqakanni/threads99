# Build stage
FROM golang:1.19-alpine AS builder

# Add git for potential private dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Build the application
WORKDIR /app/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/cmd/server/main .

# Copy config files
COPY --from=builder /app/config ./config

# Expose the port
EXPOSE 8090

# Run the binary
CMD ["./main"]