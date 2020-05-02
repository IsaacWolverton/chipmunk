FROM ubuntu:16.04 AS builder
RUN apt-get update -qq && apt-get install -qq -y curl 
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl

FROM ubuntu:16.04

COPY --from=builder /kubectl /
COPY pod_script.sh node_script.sh /
RUN chmod +x /kubectl /pod_script.sh /node_script.sh

CMD cp /node_script.sh /host; chroot /host ./node_script.sh; ./pod_script.sh