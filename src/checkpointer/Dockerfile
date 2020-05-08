# Build checkpointer
FROM golang:1.14.2 AS builder
ENV GO111MODULE=on

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download

COPY *.go ./
RUN env GOOS=linux GOARCH=amd64 go build checkpointer

FROM ubuntu:16.04
RUN apt-get update -qq && apt-get install -qq -y \
    ca-certificates \
    openssl \
    curl \
    lsb-release \
    python3

RUN export GCSFUSE_REPO=gcsfuse-`lsb_release -c -s` ;\
    echo "deb http://packages.cloud.google.com/apt $GCSFUSE_REPO main" | tee /etc/apt/sources.list.d/gcsfuse.list ;\
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
RUN apt-get update -qq && apt-get install -y -qq gcsfuse

COPY --from=builder src/checkpointer /bin/checkpointer

CMD [ "checkpointer" ]
