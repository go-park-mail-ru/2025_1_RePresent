# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем git (необходим для go mod)
RUN apk add --no-cache git

WORKDIR /workspace

# 1. Копируем ТОЛЬКО go.mod и go.sum для кэширования
COPY auth-service/go.mod auth-service/go.sum ./auth-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# 2. Устанавливаем replace для локального модуля (критично!)
RUN cd auth-service && \
    go mod edit -replace auth-service/pkg=../pkg && \
    go mod download -x

# Копируем pkg и auth-service на один уровень
COPY ./pkg ./pkg
COPY ./auth-service ./auth-service
COPY ./configs ./auth-service/configs

WORKDIR /workspace/auth-service

# 4. Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -ldflags="-w -s" -o /app/main /workspace/auth-service

# Этап запуска
FROM alpine:latest

COPY --from=builder /app/main /app/main

COPY --from=builder /workspace/auth-service/configs /app/configs

RUN chmod +x /app/main

WORKDIR /app

EXPOSE 8025

CMD ["/app/main"]