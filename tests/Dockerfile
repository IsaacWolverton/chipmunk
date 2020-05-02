# Build criu
FROM debian AS criu_builder
RUN apt-get update -qq && apt-get install -qq -y \
    gcc-multilib \
    build-essential \
    gcc \
    ccache \
    git-core \
    protobuf-c-compiler \
    protobuf-compiler \
    pkg-config 

COPY scripts/criu.package.list package.list
RUN apt-get update -qq && apt-get install -qq -y $(grep -vE "^\s*#" package.list  | tr "\n" " ")

RUN git clone https://github.com/checkpoint-restore/criu.git criu
WORKDIR /criu
RUN make


# Build chipmunk
FROM golang:1.14.2 AS chipmunk_builder
ENV GO111MODULE=on

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download

COPY main.go .
RUN env GOOS=linux GOARCH=amd64 go build chipmunk


# Main docker file
FROM debian
COPY scripts/criu.package.list package.list
RUN apt-get update -qq && apt-get install -qq -y $(grep -vE "^\s*#" package.list  | tr "\n" " ")

COPY scripts/docker.package.list package.list
RUN apt-get update -qq && apt-get install -qq -y $(grep -vE "^\s*#" package.list  | tr "\n" " ")
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN apt-key fingerprint 0EBFCD88

RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

RUN apt-get update -qq && apt-get install -qq -y docker-ce docker-ce-cli
RUN mkdir /etc/docker
RUN echo "{\"experimental\": true}" >> /etc/docker/daemon.json

COPY --from=criu_builder criu/criu/criu /bin/criu
COPY --from=chipmunk_builder src/chipmunk /bin/chipmunk


# TODO: will be the application managing checkpoints and moby will be running in the background
CMD dockerd