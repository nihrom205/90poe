# Стадия сборки
FROM golang:1.24-alpine as builder

WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN GOOS=linux go build -o app_port ./cmd/main.go

# Стадия выполнения
FROM alpine:latest

# Очень важно регулярно обновлять пакеты в изображении, чтобы включить исправления безопасности
# Также установите Bash, чтобы запустить wate-for-it.sh
RUN apk update && apk upgrade && apk add bash

#Reduce image size
RUN rm -rf /var/cache/apk/* /tmp/*

# Избегайте запуска кода в качестве пользователя root
RUN adduser -D appuser
USER appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарный файл и директорию с миграциями
COPY --from=builder /app/app_port .
COPY --from=builder /app/cmd/wait-for-it.sh .
COPY --from=builder /app/internal/app/migrations ./migrations

# Запускаем сервис
CMD ["./app_port"]