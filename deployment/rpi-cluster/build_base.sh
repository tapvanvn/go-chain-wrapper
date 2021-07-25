# we assume that 

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd "$DIR/../../"

server_url=ghcr.io/tapvanvn/go-chain-wrapper

docker build -t $server_url/rpi_chain_wrapper_base:latest -f docker/rpi_base.dockerfile ./

docker push $server_url/rpi_chain_wrapper_base:latest

popd