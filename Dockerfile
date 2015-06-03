FROM golang:1.4.2-wheezy

RUN go get github.com/lizdeika/picfit
RUN cd /go/src/github.com/lizdeika/picfit && make build

CMD /go/src/github.com/lizdeika/picfit/bin/picfit -c /etc/picfit/config.json

# EXPOSE 8080
