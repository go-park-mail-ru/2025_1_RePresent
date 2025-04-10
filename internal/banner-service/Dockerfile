# Этап сборки
FROM golang:1.24.0-alpine AS builder

# Устанавливаем git (необходим для go mod)
RUN apk add --no-cache git

WORKDIR /bannerworkspace

# 1. Копируем ТОЛЬКО go.mod и go.sum для кэширования
COPY banner-service/go.mod banner-service/go.sum ./banner-service/
COPY pkg/go.mod pkg/go.sum ./pkg/

# 2. Устанавливаем replace для локального модуля (критично!)
RUN cd banner-service && \
    go mod edit -replace pkg=../pkg && \
    go mod download -x

# Копируем pkg и auth-service на один уровень
COPY ./pkg ./pkg
COPY ./banner-service ./banner-service
COPY ./configs ./banner-service/configs

WORKDIR /bannerworkspace/banner-service

# 4. Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -ldflags="-w -s" -o /bannerapp/main /bannerworkspace/banner-service

# Этап запуска
FROM alpine:latest

COPY --from=builder /bannerapp/main /bannerapp/main

COPY --from=builder /bannerworkspace/banner-service/templates /bannerapp/templates

COPY --from=builder /bannerworkspace/banner-service/configs /bannerapp/configs

RUN chmod +x /bannerapp/main

WORKDIR /bannerapp

EXPOSE 8024

CMD ["/bannerapp/main"]