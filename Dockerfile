FROM golang:latest

# Установка gcc
RUN apt-get update && apt-get install -y gcc

# Копируем исходный код в контейнер
COPY . /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Собираем приложение
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o ./bin/linux_amd64/api ./startserv

# Указываем команду для запуска приложения
CMD ["./bin/linux_amd64/api"]