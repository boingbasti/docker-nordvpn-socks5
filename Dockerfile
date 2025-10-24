# Build stage
# Verwende die neueste stabile Go-Version (aktuell 1.25.2)
FROM golang:1.25.2-alpine AS builder
WORKDIR /app

COPY main.go .
RUN go mod init socks5 \
 && go mod tidy \
 && go build -o socks5proxy main.go

# Runtime stage
# Pinne die Alpine-Version (aktuell 3.20)
FROM alpine:3.20
WORKDIR /app/

# --- HINZUGEFÜGT ---
# Installiere curl, benötigt für den Healthcheck
RUN apk add --no-cache curl
# --- ENDE ---

# Erstelle einen dedizierten User ohne Shell
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Kopiere die Binary
COPY --from=builder /app/socks5proxy .

# Setze die Berechtigungen für den neuen User
RUN chown appuser:appgroup socks5proxy

# Wechsle zum neuen User
USER appuser

EXPOSE 1080
CMD ["./socks5proxy"]
