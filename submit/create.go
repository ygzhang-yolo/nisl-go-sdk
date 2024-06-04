package submit

import (
	"fmt"
	"nisl-go-sdk/connect"
	"nisl-go-sdk/routepool"
	"nisl-go-sdk/transactions"
	"strconv"
	"sync"
	"time"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/11 14:11
 * @Func:
 **/

func CreateLedger(ccfunc string, txNums int) {
	var wg sync.WaitGroup
	wg.Add(txNums)
	// 创建交易列表
	txs := make([]transactions.Transaction, txNums)
	for i := range txs {
		txs[i] = transactions.Transaction{
			Id:         i,
			Status:     transactions.NOTSUBMIT,
			Key:        string(i),
			Value:      "",
			CreateTime: time.Time{},
			FinishTime: time.Time{},
		}
		// ccargs
		var ccargs []string
		//ccargs = append(ccargs, []string{"CAR" + strconv.Itoa(i), "VW", "Polo", "Grey", "Mary"}...)
		ccargs = append(ccargs, []string{"owner" + strconv.Itoa(i), "10000000"}...)

		//task
		tk := &taskArgs{
			Tx:     &txs[i],
			Ct:     connect.Contract,
			Ccfunc: ccfunc,
			CCargs: ccargs,
			Wg:     &wg,
		}
		routepool.Rpool.Submit(taskWarpper(tk))
	}
	wg.Wait()
	fmt.Println("打印状态信息: ", ccfunc)
	sendrate := transactions.Transaction.GetSendRate(transactions.Transaction{}, txs)
	fmt.Printf("SendRate: %f\n", sendrate)
	tps := transactions.Transaction.GetThroughput(transactions.Transaction{}, txs)
	fmt.Printf("TPS: %f\n", tps)
	latency := transactions.Transaction.GetAvglatency(transactions.Transaction{}, txs)
	fmt.Printf("Latency: %f\n", latency)
}
