FROM alpine
MAINTAINER Cedric Lamoriniere <cedric.lamoriniere@gmail.com>

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

EXPOSE     8080
WORKDIR    /

ADD ./meraki-exporter /
ENTRYPOINT [ "/meraki-exporter" ]
CMD ["--api-freq=15s"]


