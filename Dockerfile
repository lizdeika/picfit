FROM ubuntu:14.04
MAINTAINER lizdeika

WORKDIR /opt/go/src/github.com/lizdeika/picfit
ENV GOPATH /opt/go

RUN apt-get update && apt-get install -qy \
    build-essential \
    git \
    golang

ADD . /opt/go/src/github.com/lizdeika/picfit
RUN cd /opt/go/src/github.com/lizdeika/picfit && make deps && make build

ENTRYPOINT ["/opt/go/src/github.com/lizdeika/picfit/bin/picfit"]

EXPOSE 8080
