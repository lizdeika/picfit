FROM golang

RUN apt-get update && apt-get install -qy redis-server

COPY redis.conf /usr/local/etc/redis/redis.conf
RUN /usr/bin/redis-server /usr/local/etc/redis/redis.conf

RUN mkdir /tmp/picfit

ADD . /go/src/github.com/lizdeika/picfit
RUN go get github.com/lizdeika/picfit
RUN cd /go/src/github.com/lizdeika/picfit && make build

ENTRYPOINT ["/go/src/github.com/lizdeika/picfit/bin/picfit"]

# EXPOSE 8080
