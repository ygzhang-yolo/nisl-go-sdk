# Nisl-go-sdk
自己实现了一个适用于fabric2.4的go-sdk实现, 用于论文实验;

通过线程池的方式并发提交交易, 支持设置发送速率; 统计交易的性能指标: Latency/TPS

## Quick Start

```bash
# 启动fabric网络
./startFabric.sh

# 编译项目
go mod tidy
go build

# 运行
go run nisl-go-sdk -opt=create -cck=owner0 -ccv=10  #创建一个账户<owner10, 10>
go run nisl-go-sdk -sr=100 -tx=1000  # 以100的发送速率提交1000个交易
go run nisl-go-sdk -sr=100 -st=10    # 以100的发送速率提交10s的交易

# 关闭fabric网络
./stopFabric.sh
```

## 项目说明
项目代码相关
- `connect`: 用于与区块链网络通过sdk进行连接, 需要修改其中的ccp, cca路径等信息, **里面的绝对路径要替换为本地网络的路径**;
- `logAnalyze`: 提供了一些方法, 用于分析fabric网络中的docker日志, 用于计算交易的EOV各个阶段的时延
- `routepool`: 通过ants实现的线程池来并发提交交易
- `submit`: 具体的fabric提交交易的任务
    - `create`: 创建账户的操作, 会在网络初始化时执行一次, 初始化10000个账户;
    - `invoke`: invoke操作账户的操作, 用于执行一些读写操作, 用于测试网络性能
- `transactions`: 区块链中交易的定义; 交易包含了id, status, create time, finish time; 支持计算latency, tps等指标, 具体计算方法参见其中代码部分

项目代码无关的配置文件
- `keystore`: 包含了ca信息, 用于连接区块链网络;
- `wallet`: 用于连接区块链网络的钱包信息
- `data`: 账本文件的本地存储, 区块链提交的所有交易, 均有一份记录到本地文件中以供查询;
    - `ledger.data`: 状态数据库, 存储了区块链中提交的所有kv
    - `txs`: 记录了所有交易通过sdk提交的状态, 每一条交易的提交信息和是否成功的状态信息
- `dataset`: 提供了一些数据集, 包括一些符合zipfian分布的kv数据集, 用测试不同读写集和冲突率

其他一些项目布置的脚本
- startFabric.sh: 启动fabric网络, 需要将fabric-samples目录放置在与nisl-go-sdk目录同级的位置, 否则需要自行修改对应的路径
- stopFabric.sh: 关闭fabric网络, 并删除一些账本文件和残留的文件数据
- changeConfig.sh: 用于动态更新区块链的区块配置, 如blockSize和blockNum等
