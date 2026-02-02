FROM golang:1.25 AS builder

ARG CGO_ENABLED=0
WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM golang:1.25-alpine

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /go/bin/air /usr/local/bin/

# CMD ["./main"]
CMD ["air", "-c", ".air.toml"]