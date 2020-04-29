FROM debian AS criu
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

COPY --from=criu criu/criu/criu /bin/criu

# TODO: will be the application managing checkpoints and moby will be running in the background
ENTRYPOINT [ "dockerd" ]
