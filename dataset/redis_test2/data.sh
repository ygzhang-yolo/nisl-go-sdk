#!/bin/bash

# 获取命令行输入的参数, $2是数据集大小，$3是skew参数，$4是输出文件名
args="$2 $3 $4"

# 定义脚本文件名
script="zipfian.py"

# 获取命令行输入的参数
num_runs="$1"

# 循环运行脚本指定次数
for ((i=0; i<num_runs; i++))
do
    # 传递适当的参数给脚本
    python3 $script $args
done
