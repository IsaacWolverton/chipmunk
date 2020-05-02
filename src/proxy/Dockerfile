# # Build proxt
# FROM golang:1.14.2 AS builder
# ENV GO111MODULE=on

# # COPY go.mod go.sum /src/
# COPY go.mod /src/
# WORKDIR /src
# RUN go mod download

# COPY main.go .
# RUN env GOOS=linux GOARCH=amd64 go build proxy

# FROM ubuntu:16.04
# RUN apt-get update -qq && apt-get install -qq -y \
#     ca-certificates \
#     openssl

# COPY --from=builder src/proxy /bin/proxy
# CMD [ "proxy" ]

# TODO: actually log packages, using this one for debugging
FROM golang:1.14.2 
ENV GO111MODULE=on
RUN git clone https://github.com/kahlys/proxy.git proxy
WORKDIR proxy 
RUN go build cmd/tcpproxy/*.go
CMD ["./main", "-lhost", ":42069", "-rhost", ":8080"]