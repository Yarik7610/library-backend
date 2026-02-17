FROM golang:1.25 AS builder

ARG CGO_ENABLED=0
WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/service-main ./cmd/notification-service

FROM golang:1.25-alpine

WORKDIR /app
COPY --from=builder /usr/local/bin/service-main /usr/local/bin/service-main
COPY --from=builder /go/bin/air /usr/local/bin/

# CMD ["service-main"]
CMD ["air", "-c", ".air.toml"]