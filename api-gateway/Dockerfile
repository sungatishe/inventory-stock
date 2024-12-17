FROM golang:1.22.7-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости и устанавливаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка исполняемого файла
RUN go build -o /api-gateway ./cmd/api/main.go

# Создаем минимальный образ для запуска
FROM alpine:latest
WORKDIR /root/

# Копируем скомпилированное приложение из builder'а
COPY --from=builder /api-gateway .

# Экспорт порта
EXPOSE 8080

# Команда для запуска приложения
CMD ["./api-gateway"]
