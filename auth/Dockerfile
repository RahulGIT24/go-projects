FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

COPY .env .env

RUN go build -o auth

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth .
COPY --from=builder /app/.env .env

EXPOSE 8000

CMD ["./auth"]
