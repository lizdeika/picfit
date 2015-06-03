FROM golang:1.4.2-wheezy

RUN mkdir /tmp/picfit

RUN go get github.com/lizdeika/picfit
RUN cd /go/src/github.com/lizdeika/picfit && make build

ENTRYPOINT ["/go/src/github.com/lizdeika/picfit/bin/picfit"]

# EXPOSE 8080
