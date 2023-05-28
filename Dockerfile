FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o service cmd/main/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8081

ENV PORT 8081
ENV EMAIL marvelous.mail.sender@gmail.com
ENV EMAIL_PASSWORD "twmdbnjitcszfaug"
ENV FROM_CURRENCY UTC
ENV TO_CURRENCY UAH
ENV DB_FILE_PATH ./resources/emails.dat

RUN go build

CMD ["./service"]