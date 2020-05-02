echo "starting configurator [node]"

if [ ! -f /i_was_here ]; then
    echo " -> configurating docker [node]"
    echo "{\"experimental\": true}" >> /etc/docker/daemon.json
    touch /i_was_here
    systemctl restart docker
elif [ ! -f /bin/criu ]; then
    echo " -> configurating criu [node]"

    apt-get update -qq && apt-get install -y -qq \
        gcc-multilib \
        build-essential \
        gcc \
        ccache \
        git-core \
        protobuf-c-compiler \
        protobuf-compiler \
        pkg-config \
        libnet-dev \
        libnl-route-3-dev  \
        libaio-dev \
        libcap-dev \
        libgnutls28-dev \
        libgnutls30 \
        libnl-3-dev \
        libprotobuf-c-dev \
        libprotobuf-dev \
        libselinux-dev \
        libbsd-dev \
        bsdmainutils \
        iptables \
        python-minimal \
        python-future
    
    if [ ! -d /home/root/criu ]; then
        git clone https://github.com/checkpoint-restore/criu.git /home/root/criu
    fi

    cd /home/root/criu
    make
    cp criu/criu /bin/
fi

echo "configuratation done [node]"
