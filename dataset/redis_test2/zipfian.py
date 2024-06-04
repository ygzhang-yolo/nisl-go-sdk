import sys
from scipy.stats import zipfian

# 从命令行获取参数
n = int(sys.argv[1])  # 数据集大小
a = float(sys.argv[2])  # 参数 a
filename = sys.argv[3]  # 获取命令行输入的文件名

num = zipfian.rvs(a, 10000, loc=0, size=n)

with open(filename, 'a') as file:
    for i in range(0, len(num), 1):
        item1 = num[i]
        file.write(f'{item1}\n')
