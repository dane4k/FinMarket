FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o /app/main ./app/cmd/main

CMD ["/app/main"]

LABEL authors="danya"

EXPOSE 8080