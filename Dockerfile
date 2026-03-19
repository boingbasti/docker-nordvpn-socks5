# Build stage
FROM golang:1.26.1-alpine3.23 AS builder
WORKDIR /app

COPY main.go .
RUN go mod init socks5 \
 && go mod tidy \
 && go build -o socks5proxy main.go

# Runtime stage
FROM alpine:3.23
WORKDIR /app/

# Install curl for healthcheck
RUN apk add --no-cache curl

# Create dedicated user without shell
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary
COPY --from=builder /app/socks5proxy .

# Set permissions for the new user
RUN chown appuser:appgroup socks5proxy

# Switch to non-root user
USER appuser

EXPOSE 1080
CMD ["./socks5proxy"]
