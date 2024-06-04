package submit

import (
	"fmt"
	"nisl-go-sdk/transactions"
	"sync"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/11 14:49
 * @Func:
 **/

// 使用闭包来对task进行封装
type taskFunc func()

// 执行一个task需要的参数args
type taskArgs struct {
	Tx     *transactions.Transaction //交易
	Ct     *gateway.Contract         //合约
	Ccfunc string                    //链码名
	CCargs []string                  //链码参数
	Wg     *sync.WaitGroup
}

// https://geekdaxue.co/read/startisan@go-daily-lib/crbxqg#xls1s
// 包装task func的函数
func taskWarpper(args *taskArgs) taskFunc {
	return func() {
		tx := args.Tx
		ct := args.Ct
		tx.CreateTime = time.Now()
		_, err := ct.SubmitTransaction(args.Ccfunc, args.CCargs...)
		tx.FinishTime = time.Now()
		if err != nil {
			fmt.Printf("Failed to submit transaction [id=%v, key=%v, err=%v]\n", tx.Id, tx.Key, err)
			tx.Status = transactions.FAIL
		} else {
			fmt.Printf("Success to submit transaction [id=%v, func=%v, key=%v]\n", tx.Id, args.Ccfunc, tx.Key)
			tx.Status = transactions.SUCCESS
		}
		args.Wg.Done()
	}
}
