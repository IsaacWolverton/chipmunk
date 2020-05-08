FROM ubuntu:16.04
COPY script.sh script.sh
RUN chmod +x script.sh
CMD cp script.sh /host; chroot host ./script.sh