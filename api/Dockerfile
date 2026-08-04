# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Устанавливаем необходимые утилиты сборки
RUN apk add --no-cache git ca-certificates tzdata

# Кэшируем зависимости Go
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник с отключенным CGO (статическая бинарная сборка)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api-server ./cmd/api

# --- Stage 2: Minimal Runtime Image ---
FROM alpine:3.21

WORKDIR /app

# Копируем сертификаты и часовой пояс из builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/api-server .

EXPOSE 8080

CMD ["./api-server"]
