#!/bin/bash

# 初始化计数器
counter=1

# 遍历匹配的文件
for file in 1randomtx*; do
    # 使用正则表达式提取文件名中的数字和后缀
    if [[ $file =~ 1randomtx([0-9]+)\.(.*) ]]; then
        # 构造新文件名，保留原始数字和后缀，但将前缀替换为递增的数字
        new_file="${counter}.${BASH_REMATCH[1]}.${BASH_REMATCH[2]}"
        # 重命名文件
        mv "$file" "$new_file"
        echo "Renamed $file to $new_file"
        # 更新计数器
        ((counter++))
    fi
done
