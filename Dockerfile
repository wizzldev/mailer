FROM golang:1.23.6-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/mailer

CMD ["./bin/mailer", "--listenAddr", ":7000"]
