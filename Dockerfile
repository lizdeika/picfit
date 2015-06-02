FROM golang:1.4.2-wheezy

RUN mkdir /tmp/picfit

RUN apt-get update && apt-get install -qy redis-server

ADD .  /go/src/github.com/lizdeika/picfit
RUN cd /go/src/github.com/lizdeika/picfit && make build

COPY redis.conf /usr/local/etc/redis/redis.conf
RUN /usr/bin/redis-server /usr/local/etc/redis/redis.conf

ENTRYPOINT ["/go/src/github.com/lizdeika/picfit/bin/picfit"]

# EXPOSE 8080
