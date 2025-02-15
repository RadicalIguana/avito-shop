# Этап 1: Сборка и установка зависимостей
FROM golang:1.22-alpine AS builder

# Устанавливаем зависимости и golang-migrate
RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем модули и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код приложения
COPY . .

# Переходим в папку cmd для работы с приложением
WORKDIR /app/cmd

# Этап 2: Финальный образ
FROM golang:1.22-alpine

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код и зависимости из builder
COPY --from=builder /app .

# Копируем бинарник golang-migrate
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Открываем порт 8080 для сервиса
EXPOSE 8080

# Запускаем приложение через go run
CMD ["go", "run", "cmd/main.go"]