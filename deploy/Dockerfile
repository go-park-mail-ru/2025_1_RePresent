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
# COPY pkg/templates ./internal/banner-service/templates
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/banner-service ./cmd/banner-service

# ---- Profile Service ----
FROM deps AS profile-builder
COPY cmd/profile-service ./cmd/profile-service
COPY internal/profile-service ./internal/profile-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/profile-service ./cmd/profile-service

# ---- Pay Service ----
FROM deps AS pay-builder
COPY cmd/pay-service ./cmd/pay-service
COPY internal/pay-service ./internal/pay-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/pay-service ./cmd/pay-service

# ---- CSAT Service ----
FROM deps AS csat-builder
COPY cmd/csat-service ./cmd/csat-service
COPY internal/csat-service ./internal/csat-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/csat-service ./cmd/csat-service

# ---- Adv Service ----
FROM deps AS adv-builder
COPY cmd/adv-service ./cmd/adv-service
COPY internal/adv-service ./internal/adv-service
COPY pkg ./pkg
COPY configs ./configs
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/adv-service ./cmd/adv-service

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
COPY --from=banner-builder /app/pkg/iframe-templates ./templates
ENTRYPOINT ["./banner-service"]

# ---- Profile Service ----
FROM alpine:latest AS profile-service
WORKDIR /app
COPY --from=profile-builder /bin/profile-service ./
COPY --from=profile-builder /app/configs ./configs
ENTRYPOINT ["./profile-service"]

# ---- Pay Service ----
FROM alpine:latest AS pay-service
WORKDIR /app
COPY --from=pay-builder /bin/pay-service ./
COPY --from=pay-builder /app/configs ./configs
ENTRYPOINT ["./pay-service"]

# ---- CSAT Service ----
FROM alpine:latest AS csat-service
WORKDIR /app
COPY --from=csat-builder /bin/csat-service ./
COPY --from=csat-builder /app/configs ./configs
COPY --from=csat-builder /app/internal/csat-service/templates ./templates
ENTRYPOINT ["./csat-service"]

# ---- Adv Service ----
FROM alpine:latest AS adv-service
WORKDIR /app
COPY --from=adv-builder /bin/adv-service ./
COPY --from=adv-builder /app/configs ./configs
COPY --from=banner-builder /app/pkg/iframe-templates ./templates
ENTRYPOINT ["./adv-service"]
