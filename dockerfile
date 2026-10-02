# --- Этап сборки ---
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Сначала копируем манифесты — так кэш слоёв работает эффективнее
COPY go.mod go.sum ./
RUN go mod download

# Затем исходники
COPY . .

# Собираем статически слинкованный бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/daysny .

# --- Финальный образ ---
FROM alpine:3.20

RUN adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=builder /out/daysny /app/daysny

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/daysny"]