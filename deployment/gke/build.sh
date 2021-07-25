# we assume that 

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd "$DIR/../../"

tag=$(<./version.txt)

server_url=ghcr.io/tapvanvn/go-chain-wrapper

docker build -t $server_url/gke_chain_wrapper:$tag -f docker/gke.dockerfile ./

docker push $server_url/gke_chain_wrapper:$tag

popd