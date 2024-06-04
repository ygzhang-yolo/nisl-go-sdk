#!/bin/bash

# 定义脚本文件名的前缀和后缀
script_prefix="./data.sh"
script_suffix="randomtx"

# 定义参数列表
param_values=("0.2" "0.4" "0.6" "0.8" "1.0" "1.2" "1.4" "1.6" "1.8" "2.0")

# 定义函数以记录脚本运行时间
record_script_time() {
    script_name=$1
    start_time=$(date +%s.%N)
    # 运行脚本
    ./$script_name
    end_time=$(date +%s.%N)
    # 计算运行时间
    duration=$(echo "$end_time - $start_time" | bc)
    # 输出结果
    echo "脚本 $script_name 运行时间: $duration 秒"
}

# 遍历参数列表，生成并运行脚本
for param_value in "${param_values[@]}"; do
    script_name="${script_prefix} 200 100 ${param_value} ${script_suffix}${param_value}"
    record_script_time "$script_name" &
done

# 等待所有后台任务完成
wait
