FROM golang:buster AS builder

RUN go version

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates git zlib1g-dev

COPY . /go/src/github.com/TicketsBot/data-removal-service
WORKDIR  /go/src/github.com/TicketsBot/data-removal-service

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -o service cmd/data-removal-service/main.go

# Prod container
FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates curl

COPY --from=builder /go/src/github.com/TicketsBot/data-removal-service/service /srv/data-removal-service/service
RUN chmod +x /srv/data-removal-service/service

RUN useradd -m container
USER container
WORKDIR /srv/data-removal-service

CMD ["/srv/data-removal-service/service"]