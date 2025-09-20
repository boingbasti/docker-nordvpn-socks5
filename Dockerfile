# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY main.go .
RUN go mod init socks5 \
 && go mod tidy \
 && go build -o socks5proxy main.go

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/socks5proxy .
EXPOSE 1080
CMD ["./socks5proxy"]
