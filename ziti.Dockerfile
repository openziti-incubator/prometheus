ARG ARCH="amd64"
ARG OS="linux"

FROM golang as build
RUN curl -sL https://deb.nodesource.com/setup_16.x | bash -
RUN apt-get install -y bzip2 nodejs
ADD . /build
WORKDIR /build
# openziti:foundation requires cgo.
# enable cgo by hijacking the PROMU_BINARIES variable to pass `--cgo` to promu
RUN rm -rf web/ui/node_modules
RUN make build PROMU_BINARIES=--cgo


###
# The final stage is a slightly modified version of
#   https://github.com/openziti-incubator/prometheus/blob/main/Dockerfile
#

FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="NetFoundry - The Prometheus Zitifiers <adv-dev@netfoundry.io>"

ARG ARCH="amd64"
ARG OS="linux"
COPY --from=build /build/prometheus                             /bin/prometheus
COPY --from=build /build/promtool                               /bin/promtool
COPY --from=build /build/documentation/examples/prometheus.yml  /etc/prometheus/prometheus.yml
COPY --from=build /build/console_libraries/                     /usr/share/prometheus/console_libraries/
COPY --from=build /build/consoles/                              /usr/share/prometheus/consoles/
COPY --from=build /build/LICENSE                                /LICENSE
COPY --from=build /build/NOTICE                                 /NOTICE
COPY --from=build /build/npm_licenses.tar.bz2                   /npm_licenses.tar.bz2

WORKDIR /prometheus
RUN ln -s /usr/share/prometheus/console_libraries /usr/share/prometheus/consoles/ /etc/prometheus/ && \
    chown -R nobody:nobody /etc/prometheus /prometheus

USER       nobody
EXPOSE     9090
VOLUME     [ "/prometheus" ]
ENTRYPOINT [ "/bin/prometheus" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles" ]
