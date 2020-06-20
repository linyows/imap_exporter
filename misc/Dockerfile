ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc

ARG ARCH="amd64"
ARG OS="linux"
COPY ./.build/${OS}-${ARCH}/imap_exporter /bin/imap_exporter
COPY ./config.yml /imap_exporter.yml

EXPOSE      9993
USER        nobody
ENTRYPOINT  [ "/bin/imap_exporter", "--config", "/imap_exporter.yml" ]
CMD [ "--imap-address", "pop3.lolipop.jp:993", "--log.level", "debug" ]
