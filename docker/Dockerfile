FROM golang:1.12-alpine as build
MAINTAINER dev <aet> fedux dot org

RUN apk update \
    && apk upgrade \
    && apk add git gcc make g++

COPY . /build/local-webserver
WORKDIR /build/local-webserver
RUN bin/build

####

FROM alpine:3.11

RUN apk update --no-cache \
    && apk upgrade --no-cache \
    && apk add tini

COPY --from=build /build/local-webserver/dist/linux/amd64/lw /usr/local/bin/lw
RUN chmod +x /usr/local/bin/lw

EXPOSE 80

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/local/bin/lw", "--no-browser", "--interface=0.0.0.0", "--port=80", "--directory=/data"]
