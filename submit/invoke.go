package submit

import (
	"context"
	"fmt"
	"nisl-go-sdk/connect"
	"nisl-go-sdk/dataset"
	"nisl-go-sdk/routepool"
	"nisl-go-sdk/transactions"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/11 15:41
 * @Func:
 **/

const DataPath = "data/txs/"
const TempDataPath = "data/temp/"
const DatasetPath = "dataset/testdata/"

func InvokeTx(txNums int, ccfunc string) {
	//txNums := 10
	var wg sync.WaitGroup
	wg.Add(txNums)
	// 创建交易列表
	txs := make([]transactions.Transaction, txNums)
	for i := range txs {
		txs[i] = transactions.Transaction{
			Id:         i,
			Status:     transactions.NOTSUBMIT,
			Key:        string(i),
			CreateTime: time.Time{},
			FinishTime: time.Time{},
		}
		// ccargs
		var ccargs []string
		ccargs = append(ccargs, []string{"CAR" + strconv.Itoa(i), "Archie"}...)

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

func invokeTxWithSR(txs []transactions.Transaction, ccfunc string, sr int) {
	//控制发送速率
	limiter := rate.NewLimiter(rate.Limit(sr), 1)
	//设置上下文，超时的时间为100s
	ctx, _ := context.WithTimeout(context.Background(), time.Second*100)
	txNums := len(txs)
	var wg sync.WaitGroup
	wg.Add(txNums)
	// 创建交易列表
	for i := range txs {
		// ccargs
		var ccargs []string
		ccargs = append(ccargs, []string{"CAR" + strconv.Itoa(i), "Archie"}...)
		//task
		tk := &taskArgs{
			Tx:     &txs[i],
			Ct:     connect.Contract,
			Ccfunc: ccfunc,
			CCargs: ccargs,
			Wg:     &wg,
		}
		// 等待时间间隔，控制发送速率
		limiter.WaitN(ctx, 1)
		routepool.Rpool.Submit(taskWarpper(tk))
	}
	wg.Wait()
}

func invokeTxWithSRCR(txs []transactions.Transaction, sr int) {
	//控制发送速率
	limiter := rate.NewLimiter(rate.Limit(sr), 1)
	//设置上下文，超时的时间为100s
	ctx, _ := context.WithTimeout(context.Background(), time.Second*100)
	txNums := len(txs)
	var wg sync.WaitGroup
	wg.Add(txNums)
	// 创建交易列表
	for i := range txs {
		// ccargs & ccfunc
		var ccfunc string
		var ccargs []string
		keys := strings.Split(txs[i].Key, ",")
		if txs[i].CCFunc == "0" {
			ccfunc = "SEND_PAYMENT"
			ccargs = append(ccargs, []string{keys[0], keys[1], "1"}...)
		} else {
			ccfunc = "BALANCE"
			ccargs = append(ccargs, []string{keys[0]}...)
		}
		//task
		tk := &taskArgs{
			Tx:     &txs[i],
			Ct:     connect.Contract,
			Ccfunc: ccfunc,
			CCargs: ccargs,
			Wg:     &wg,
		}
		// 等待时间间隔，控制发送速率
		limiter.WaitN(ctx, 1)
		routepool.Rpool.Submit(taskWarpper(tk))
	}
	wg.Wait()
}

func SimpleSubmitTxWithSR(txNums int, ccfunc string, sr int) {
	// 创建交易列表
	txs := make([]transactions.Transaction, txNums)
	for i := range txs {
		txs[i] = transactions.Transaction{
			Id:         i,
			Status:     transactions.NOTSUBMIT,
			Key:        string(i),
			CreateTime: time.Time{},
			FinishTime: time.Time{},
		}
	}
	invokeTxWithSR(txs, ccfunc, sr)
}

func SubmitZipfianTxWithSR(txNums int, ccfunc string, sr int) {
	keys := dataset.GetDataFromFile("/home/zhangyiguang/fabric/nisl-go-sdk/dataset/zipfian.data")
	// 先简单选取200个
	keys = keys[100 : 100+txNums]
	txs := make([]transactions.Transaction, len(keys))
	for i := range txs {
		txs[i] = transactions.Transaction{
			Id:         i,
			Status:     transactions.NOTSUBMIT,
			Key:        "CAR" + keys[i],
			CreateTime: time.Time{},
			FinishTime: time.Time{},
		}
	}
	invokeTxWithSR(txs, ccfunc, sr)
	fmt.Println("打印状态信息: ", ccfunc)
	sendrate := transactions.Transaction.GetSendRate(transactions.Transaction{}, txs)
	fmt.Printf("SendRate: %f\n", sendrate)
	tps := transactions.Transaction.GetThroughput(transactions.Transaction{}, txs)
	fmt.Printf("TPS: %f\n", tps)
	latency := transactions.Transaction.GetAvglatency(transactions.Transaction{}, txs)
	fmt.Printf("Latency: %f\n", latency)
	// 将结果写入文件
	writeToFile(DataPath, txs)
	writeToTempFile(TempDataPath, txs) //同时写入一个临时文件
}

func SubmitZipfianTxWithSRCR(sr int, cr string, st int) {
	keys := dataset.GetAllDataFromFile(DatasetPath + cr)
	// 控制发送时间最多60秒
	var txNums int
	if len(keys) > sr*st {
		txNums = sr * st
	} else {
		txNums = len(keys)
	}
	keys = keys[:txNums]
	txs := make([]transactions.Transaction, len(keys))
	for i, key := range keys {
		txs[i] = transactions.Transaction{
			Id:         i,
			Status:     transactions.NOTSUBMIT,
			Key:        "owner" + key[1] + "," + "owner" + key[2],
			CreateTime: time.Time{},
			FinishTime: time.Time{},
			CCFunc:     key[0],
		}
	}
	invokeTxWithSRCR(txs, sr)
	//fmt.Println("打印状态信息: ", cr)
	//sendrate := transactions.Transaction.GetSendRate(transactions.Transaction{}, txs)
	//fmt.Printf("SendRate: %f\n", sendrate)
	//tps := transactions.Transaction.GetThroughput(transactions.Transaction{}, txs)
	//fmt.Printf("TPS: %f\n", tps)
	//latency := transactions.Transaction.GetAvglatency(transactions.Transaction{}, txs)
	//fmt.Printf("Latency: %f\n", latency)
	// 将结果写入文件
	writeToFile(DataPath, txs)
	writeToTempFile(TempDataPath, txs) //同时写入一个临时文件
}

// 只创建一个smallbank的交易
func CreateSingleSmallbankTx(sr int, key string) {
	txNums := 1
	var wg sync.WaitGroup
	wg.Add(txNums)
	// 创建交易列表
	txs := make([]transactions.Transaction, txNums)
	txs[0] = transactions.Transaction{
		Id:         0,
		Status:     transactions.NOTSUBMIT,
		Key:        key,
		Value:      "",
		CreateTime: time.Time{},
		FinishTime: time.Time{},
		CCFunc:     "CreateAccount",
	}
	var ccargs []string
	ccargs = append(ccargs, []string{"owner" + key, "10000000"}...)

	//task
	tk := &taskArgs{
		Tx:     &txs[0],
		Ct:     connect.Contract,
		Ccfunc: "CreateAccount",
		CCargs: ccargs,
		Wg:     &wg,
	}
	routepool.Rpool.Submit(taskWarpper(tk))
	wg.Wait()
	// 将结果写入文件
	writeToFile(DataPath, txs)
	writeToTempFile(TempDataPath, txs) //同时写入一个临时文件
}

// 只提交一个smallbank的交易
func SubmitSingleSmallbankTx(sr int, key string, val string) {
	txNums := 1
	var wg sync.WaitGroup
	wg.Add(txNums)
	// 创建交易列表
	txs := make([]transactions.Transaction, txNums)
	txs[0] = transactions.Transaction{
		Id:         0,
		Status:     transactions.NOTSUBMIT,
		Key:        key,
		Value:      val,
		CreateTime: time.Time{},
		FinishTime: time.Time{},
	}
	var ccargs []string
	ccargs = append(ccargs, []string{"owner" + key, val}...)

	//task
	tk := &taskArgs{
		Tx:     &txs[0],
		Ct:     connect.Contract,
		Ccfunc: "TRANSACT_SAVINGS",
		CCargs: ccargs,
		Wg:     &wg,
	}
	routepool.Rpool.Submit(taskWarpper(tk))
	wg.Wait()
	// 将结果写入文件
	writeToFile(DataPath, txs)
	writeToTempFile(TempDataPath, txs) //同时写入一个临时文件
}
