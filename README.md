# go-chain-wrapper
What happen on chain 

#Run local
1. download geth from binance and locate it in ./bsc/geth
2. clone ./config/env_example.sh to ./config/env.sh
3. run ./deployment/local/run_local.sh

#Run on raspberry-pi 32bit cluster
1. change docker image repository from ./docker/rpi.dockerfile
2. change docker image repository url in ./deployment/rpi_cluster/build_base.sh && build.sh
3. run ./deployment/rpi_cluster/build_base.sh to build base image
4. run ./deployment/rpi_cluster/build.sh to build image
5. run ./deployment/rpi_cluster/helm_install.sh to deploy