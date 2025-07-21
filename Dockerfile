# 1. Используем официальный образ Go
FROM golang:1.21-alpine

# 2. Установка зависимостей
RUN apk add --no-cache git

# 3. Устанавливаем рабочую директорию
WORKDIR /app

# 4. Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 5. Копируем исходный код
COPY . .

# 6. Собираем Go-приложение
RUN go build -o server ./cmd/main.go

# 7. Создаем папку для загрузок (если нужно)
#RUN mkdir -p uploads

# 8. Экспонируем порт
EXPOSE ${PORT}

# 9. Устанавливаем переменные окружения по умолчанию
ENV PORT=8080

# 10. Запускаем сервер
CMD ["./server"]
