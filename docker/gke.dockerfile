
FROM tapvanvn/gke_jsonrpc_wrapper_base:latest AS build

WORKDIR /


RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./... && go build

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build               /src/go-jsonrpc-wrapper / 
#COPY --from=build               /3rd/bsc/build/bin/     /3rd/bsc/
COPY config/config.json        /config/config.json 
COPY config/route.json        /config/route.json 
COPY abi_file/                  /abi_file

ENTRYPOINT ["/go-jsonrpc-wrapper"]