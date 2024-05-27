FROM golang:latest

RUN apt-get update && apt-get install -y gcc

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o ./bin/linux_amd64/api ./run

CMD ["./bin/linux_amd64/api"]