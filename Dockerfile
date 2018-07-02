FROM golang:1.10

RUN mkdir -p /go/
WORKDIR /go/
COPY . ./src/github.com/waveywaves/webcrawler/

RUN export GOPATH=$(pwd)
RUN go install github.com/waveywaves/webcrawler


ENTRYPOINT ["/go/bin/webcrawler"]