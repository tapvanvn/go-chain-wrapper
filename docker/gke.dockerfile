
FROM ghcr.io/tapvanvn/go-chain-wrapper/gke_chain_wrapper_base:latest AS build

WORKDIR /


RUN mkdir -p /src

COPY ./ /src/

RUN cd /src && go get ./... && go build

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build               /src/go-chain-wrapper / 
COPY config/config.jsonc        /config/config.jsonc 
COPY config/route.jsonc        /config/route.jsonc 
COPY abi_file/                  /abi_file

ENTRYPOINT ["/go-chain-wrapper"]