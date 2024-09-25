FROM quay.io/prometheus/busybox:latest
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

COPY object-storage /bin/object-storage

ENTRYPOINT ["/bin/object-storage"]
