DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
tag=$(<../../version.txt)
name=jsonrpc_wrapper
if [ -f "$DIR/../../name.txt" ]; then
    name=$(<$DIR/../../name.txt)
fi

namespace=default
if [ -f "$DIR/../../namespace.txt" ]; then
    namespace=$(<$DIR/../../namespace.txt)
fi

helm upgrade $name $DIR/../../helm/rpi-cluster --set image.tag=$tag,image.namespace=$namespace --namespace=$namespace