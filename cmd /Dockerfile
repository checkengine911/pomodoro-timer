# Сборка Go-приложения
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY ../../internal ../internal
COPY . .
RUN go build -o server

# Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/internal /app/internal
EXPOSE 8080
CMD ["./server"]