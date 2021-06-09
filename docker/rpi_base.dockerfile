FROM arm32v7/golang:1.16-alpine

WORKDIR /

RUN apk update && apk add --no-cache git curl openssh-client gcc g++ musl-dev 
RUN apk add make linux-headers

RUN mkdir -p /3rd/bsc

RUN cd /3rd/bsc && git clone https://github.com/binance-chain/bsc ./ && make geth