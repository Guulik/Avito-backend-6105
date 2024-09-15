FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app ./cmd/tender

EXPOSE 8080

ENV SERVER_ADDRESS="0.0.0.0:8080"
ENV POSTGRES_CONN="postgres://cnrprod1725725750-team-77090:cnrprod1725725750-team-77090@postgres:6432/cnrprod1725725750-team-77090"
ENV POSTGRES_JDBC_URL="jdbc:postgresql://postgres:6432/cnrprod1725725750-team-77090"
ENV POSTGRES_USERNAME="cnrprod1725725750-team-77090"
ENV POSTGRES_PASSWORD="cnrprod1725725750-team-77090"
ENV POSTGRES_HOST="rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net"
ENV POSTGRES_PORT="6432"
ENV POSTGRES_DATABASE="cnrprod1725725750-team-77090"

CMD ["./app"]

