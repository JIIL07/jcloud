FROM golang:latest AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y gcc
RUN go build -o serve .

FROM golang:latest
RUN apt-get update && apt-get install -y ca-certificates bash
WORKDIR /root/

COPY --from=builder /app/serve .

EXPOSE 8080

CMD ["./serve"]
