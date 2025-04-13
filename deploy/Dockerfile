# ===================================
# Этап 1: Общие зависимости
# ===================================
FROM golang:1.24 AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# ===================================
# Этап 2: Сборка сервисов (каждый в своем этапе)
# ===================================

# ---- Auth Service ----
FROM deps AS auth-builder
COPY cmd/auth-service ./cmd/auth-service
COPY internal/auth-service ./internal/auth-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd/auth-service

# ---- Banner Service ----
FROM deps AS banner-builder
COPY cmd/banner-service ./cmd/banner-service
COPY internal/banner-service ./internal/banner-service
COPY pkg ./pkg
COPY configs ./configs
COPY internal/banner-service/templates ./internal/banner-service/templates
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/banner-service ./cmd/banner-service

# ---- Profile Service ----
FROM deps AS profile-builder
COPY cmd/profile-service ./cmd/profile-service
COPY internal/profile-service ./internal/profile-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/profile-service ./cmd/profile-service

# ---- Payment Service ----
FROM deps AS payment-builder
COPY cmd/payment-service ./cmd/payment-service
COPY internal/pay-service ./internal/pay-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/payment-service ./cmd/pay-service

# ===================================
# Этап 3: Финальные образы
# ===================================

# ---- Auth Service ----
FROM alpine:latest AS auth-service
WORKDIR /app
COPY --from=auth-builder /bin/auth-service ./
COPY --from=auth-builder /app/configs ./configs
ENTRYPOINT ["./auth-service"]

# ---- Banner Service ----
FROM alpine:latest AS banner-service
WORKDIR /app
COPY --from=banner-builder /bin/banner-service ./
COPY --from=banner-builder /app/configs ./configs
COPY --from=banner-builder /app/internal/banner-service/templates ./templates
ENTRYPOINT ["./banner-service"]

# ---- Profile Service ----
FROM alpine:latest AS profile-service
WORKDIR /app
COPY --from=profile-builder /bin/profile-service ./
COPY --from=profile-builder /app/configs ./configs
ENTRYPOINT ["./profile-service"]

# ---- Payment Service ----
FROM alpine:latest AS payment-service
WORKDIR /app
COPY --from=payment-builder /bin/payment-service ./
COPY --from=payment-builder /app/configs ./configs
ENTRYPOINT ["./payment-service"]