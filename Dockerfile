ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc
LABEL maintainer="Simon Schneider <dev@raynigon.com>"

COPY climate-metrics /bin/climate-metrics

EXPOSE      9776
USER        nobody
ENTRYPOINT  [ "/bin/climate-metrics" ]