FROM debian:8.7

ENV MOUNT_POINT /mountApp
RUN mkdir "${MOUNT_POINT}"
COPY main.sh /main.sh
CMD ["bash", "/main.sh"]
