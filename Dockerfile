FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app ./cmd/tender

EXPOSE 8080

CMD ["./app"]

