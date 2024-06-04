#!/bin/bash

# 遍历以"rand"开头的文件
for file in rand*; do
    # 检查文件是否存在并且是普通文件
    if [[ -f "$file" ]]; then
        line_number=0
        # 逐行读取文件内容
        while IFS= read -r line; do
            ((line_number++))
            # 使用正则表达式检查每行是否只有两个以空格相隔开的参数，并且参数为1-10000之间的正整数
            if [[ ! "$line" =~ ^([1-9][0-9]{0,3}|10000)[[:space:]]([1-9][0-9]{0,3}|10000)$ ]]; then
                echo "文件 $file 第 $line_number 行不符合要求：$line"
            fi
        done < "$file"
    fi
done
