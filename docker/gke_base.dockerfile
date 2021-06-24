FROM golang:1.16-buster AS build

WORKDIR /

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y git curl openssh-client gcc g++ musl-dev

RUN mkdir -p /3rd/bsc

RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./...

RUN rm -rf /src