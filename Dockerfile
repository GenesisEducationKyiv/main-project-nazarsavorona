FROM golang:latest

WORKDIR /api

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 8081
ENV EMAIL marvelous.mail.sender@gmail.com
ENV EMAIL_PASSWORD "twmdbnjitcszfaug"

RUN go build

CMD ["./BTCRateCheckService"]