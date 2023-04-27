# syntax=docker/dockerfile:1
FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src/ ./src

RUN go build -o fractage src/main.go

EXPOSE 6060

ENTRYPOINT ["/app/fractage"]
