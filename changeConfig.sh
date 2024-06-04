#!/bin/bash

# 判断参数数量是否正确
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <blocksize> <batch_timeout>"
    exit 1
fi

# 提取输入参数
num1=$1
num2=$2

WORKPATH=/home/zhangyiguang/fabric/nisl-go-sdk
NETWORKPATH=/home/zhangyiguang/fabric/fabric-samples/test-network

# 获取容器ID或名称
container_id=$(docker ps -qf "name=cli")

# 判断输入参数是否为不低于 0 的整数
# if ! [[ "$num1" =~ ^[0-9]+$ && "$num2" =~ ^[0-9]+$ ]]; then
#     echo "Error: Both input parameters must be non-negative integers."
#     exit 1
# fi

docker cp $WORKPATH/chanconfig.sh $container_id:/opt/gopath/src/github.com/hyperledger/fabric/peer/chanconfig.sh
docker exec -it $container_id bash -c "./chanconfig.sh $num1 $num2"