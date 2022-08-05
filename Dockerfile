FROM golang:1.10
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY . /go/src/github.com/sofyan48/mq-router
WORKDIR /go/src/github.com/sofyan48/mq-router
RUN dep ensure
