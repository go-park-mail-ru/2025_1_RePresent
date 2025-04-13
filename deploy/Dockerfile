# ===================================
# Этап 1: Общие зависимости
# ===================================
FROM golang:1.24 AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# ===================================
# Этап 2: Отдельные builder-этапы для каждого сервиса
# ===================================

# ---- Auth Service ----
FROM deps AS auth-builder
COPY cmd/auth-service ./cmd/auth-service
COPY internal/auth-service ./internal/auth-service
COPY pkg/ ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd/auth-service

# ---- Banner Service ----
FROM deps AS banner-builder
COPY cmd/banner-service ./cmd/banner-service
COPY internal/banner-service ./internal/banner-service
COPY pkg/ ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/banner-service ./cmd/banner-service

# ---- Profile Service ----
FROM deps AS profile-builder
COPY cmd/profile-service ./cmd/profile-service
COPY internal/profile-service ./internal/profile-service
COPY pkg/ ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/profile-service ./cmd/profile-service

# ---- Payment Service ----
FROM deps AS payment-builder
COPY cmd/payment-service ./cmd/payment-service
COPY internal/payment-service ./internal/payment-service
COPY pkg/ ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/payment-service ./cmd/payment-service

# ===================================
# Этап 3: Финальные образы (остаётся без изменений)
# ===================================

FROM alpine:latest AS auth-service
WORKDIR /app
COPY --from=auth-builder /bin/auth-service ./auth-service
COPY --from=auth-builder /app/configs ./configs
ENTRYPOINT ["./auth-service"]
