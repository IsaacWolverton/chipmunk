# Build checkpointer
FROM golang:1.14.2 AS builder
ENV GO111MODULE=on

# COPY go.mod go.sum /src/
COPY go.mod /src/
WORKDIR /src
RUN go mod download

COPY main.go .
RUN env GOOS=linux GOARCH=amd64 go build proxy

FROM ubuntu:16.04
RUN apt-get update -qq && apt-get install -qq -y \
    ca-certificates \
    openssl

COPY --from=builder src/proxy /bin/proxy
CMD [ "proxy" ]