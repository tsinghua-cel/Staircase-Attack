FROM ubuntu:22.04

WORKDIR /root

COPY ./bin/prysmctl /usr/bin/prysmctl
COPY ./bin/bootnode /usr/bin/bootnode

COPY ./entrypoint/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod u+x /usr/local/bin/entrypoint.sh

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
