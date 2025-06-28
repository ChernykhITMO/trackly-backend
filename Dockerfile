# Используем официальный образ Go
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

COPY . .

# генерим код
RUN oapi-codegen -generate types,server,client -package api /app/open-api/openapi.yaml > /app/internal/api/api.gen.go
# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o my-go-service ./cmd/main.go

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из builder
COPY --from=builder /app/my-go-service .

# Копируем папку с конфигурацией
COPY configs ./configs
COPY migrations ./migrations
# Открываем порт
EXPOSE 8080

# Запускаем приложение с аргументом (по умолчанию)
CMD ["./my-go-service", "--configs", "./configs/config.yaml"]