FROM busybox:latest

MAINTAINER tim@magnetic.io

ADD ./target/linux_i386/tatanka /tatanka
ADD ./config /config

ENTRYPOINT ["/tatanka"]
