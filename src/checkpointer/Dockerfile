# Build checkpointer
FROM golang:1.14.2 AS builder
ENV GO111MODULE=on

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download

COPY main.go .
RUN env GOOS=linux GOARCH=amd64 go build checkpointer

FROM ubuntu:16.04
RUN apt-get update -qq && apt-get install -qq -y \
    ca-certificates \
    openssl

COPY --from=builder src/checkpointer /bin/checkpointer
CMD [ "checkpointer" ]