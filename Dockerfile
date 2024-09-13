FROM golang:latest

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/server/main.go

RUN mkdir -p /app/secrets
COPY secrets/.env /app/secrets/.env
ENV ENV_PATH=/app/secrets/.env

RUN mkdir -p /app/config
COPY config/config.yaml /app/config/config.yaml
ENV CONFIG_PATH=/app/config/config.yaml


RUN mkdir -p /app/storage
COPY storage/storage.db /app/storage/storage.db
ENV DATABASE_PATH=/app/storage/storage.db

RUN mkdir -p /app/static
COPY cmd/server/static /app/static


EXPOSE 8080

CMD ["./server"]