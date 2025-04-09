FROM golang:1.24 AS deps

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24 AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-service ./cmd/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/banner-service ./cmd/banner-service

# ===================================
# Этап 3: Минимальные образы для сервисов
# ===================================

# ---- Service 1 ----
FROM alpine:latest AS auth-service

WORKDIR /app
COPY --from=builder /bin/auth-service ./auth-service

ENTRYPOINT ["./auth-service"]

# ---- Service 2 ----
FROM alpine:latest AS banner-service

WORKDIR /app
COPY --from=builder /bin/banner-service ./banner-service

ENTRYPOINT ["./banner-service"]