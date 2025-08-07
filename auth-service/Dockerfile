FROM golang:1.24 AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]