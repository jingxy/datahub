# This is a local-build docker image for p2p-dl test

FROM golang:1.5
MAINTAINER Zonesan <chaizs@asiainfo.com>

ENV SRCPATH $GOPATH/src/github.com/asiainfoLDP/datahub 
ENV PATH $PATH:$GOPATH/bin:$SRCPATH
RUN mkdir $SRCPATH -p
WORKDIR $SRCPATH

ADD . $SRCPATH

RUN mkdir /var/lib/datahub
RUN go build

EXPOSE 35800

CMD $SRCPATH/start.sh


