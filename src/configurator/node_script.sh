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

    echo " -> configurating gcsfuse [node]"
    
    # apt-get update -qq && apt-get install -y -qq curl
    # export GCSFUSE_REPO=gcsfuse-`lsb_release -c -s`
    # echo "deb http://packages.cloud.google.com/apt $GCSFUSE_REPO main" | tee /etc/apt/sources.list.d/gcsfuse.list
    # curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
    # apt-get update -qq && apt-get install -y -qq gcsfuse

    rm -rf /shared; mkdir /shared; chmod 777 /shared
    # gcsfuse --implicit-dirs chipmunk-storage /shared


    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
    sudo apt-get install -qq -y apt-transport-https ca-certificates gnupg
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
    sudo apt-get update -qq && sudo apt-get -qq -y install google-cloud-sdk
    gcloud config set project 'mit-mic'


fi

echo "configuratation done [node]"
