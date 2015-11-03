FROM golang:1.5.1-wheezy

ENV HONEYBADGER_ENV production

RUN go get github.com/lizdeika/picfit
RUN cd /go/src/github.com/lizdeika/picfit && make build

CMD /go/src/github.com/lizdeika/picfit/bin/picfit -c /etc/picfit/config.json

# EXPOSE 8080
