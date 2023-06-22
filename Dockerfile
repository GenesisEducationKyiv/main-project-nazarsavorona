FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o service cmd/main/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/service .
COPY --from=builder /app/templates ./templates

EXPOSE 8081

ENV PORT 8081
ENV EMAIL your.email@example.com
ENV EMAIL_PASSWORD your_password
ENV FROM_CURRENCY BTC
ENV TO_CURRENCY UAH
ENV DB_FILE_PATH ./resources/emails.dat

CMD ["./service"]