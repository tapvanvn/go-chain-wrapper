# we assume that 

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd "$DIR/../../"

server_url=tapvanvn

docker build -t $server_url/rpi_jsonrpc_wrapper_base:latest -f docker/rpi_base.dockerfile ./

docker push $server_url/rpi_jsonrpc_wrapper_base:latest

popd