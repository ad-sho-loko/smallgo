FROM ubuntu
MAINTAINER adsholoko

RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-get install -y build-essential golang-1.13-go gdb git

ENV PATH $PATH:/usr/lib/go-1.13/bin

RUN go get github.com/stretchr/testify
