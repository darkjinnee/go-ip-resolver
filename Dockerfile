# Многоэтапная сборка для Go приложения
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o ip-resolver \
    ./cmd/ip-resolver

# Финальный образ
FROM alpine:3.18

# Устанавливаем необходимые пакеты для runtime
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из builder этапа
COPY --from=builder /app/ip-resolver .

# Копируем конфигурационный файл
COPY --from=builder /app/configs/domains.json ./configs/

# Создаем директорию для логов
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Устанавливаем переменные окружения
ENV GIN_MODE=release
ENV TZ=UTC

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/groups || exit 1

# Запускаем приложение
CMD ["./ip-resolver", "-config", "configs/domains.json", "-addr", ":8080"]
