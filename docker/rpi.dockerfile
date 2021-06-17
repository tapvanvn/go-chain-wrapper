
FROM arm32v7/golang:1.16-alpine AS build

WORKDIR /

RUN apk update && apk add --no-cache git curl openssh-client gcc g++ musl-dev 
RUN apk add make linux-headers

RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./... && go build

FROM arm32v7/alpine
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

RUN mkdir -p /3rd/bsc

COPY --from=build               /src/go-jsonrpc-wrapper / 
COPY --from=build               /3rd/bsc/build/bin/     /3rd/bsc/
COPY config/config.json        /config/config.json 
COPY config/route.json        /config/route.json 
COPY abi_file/                  /abi_file

ENTRYPOINT ["/go-jsonrpc-wrapper"]