package main

import (
	"flag"
	"fmt"
	"nisl-go-sdk/connect"
	"nisl-go-sdk/routepool"
	"nisl-go-sdk/submit"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/11 14:01
 * @Func:
 **/

const MAX_ROUTINE_NUM int = 1000       //最大协程数
var ChannelName string = "mychannel"   //默认通道名
var ChaincodeName string = "smallbank" //默认链码名

var SendRate int
var ConflictRate string
var TxNums int
var SendTime int
var CCkey string //调用cc的key
var CCVal string //调用cc的val
var Opt string   //判断是否是create操作

func main() {
	// cmd args
	flag.IntVar(&SendRate, "sr", 30, "send rate")
	flag.StringVar(&ConflictRate, "cr", "init file", "conflict rate")
	flag.IntVar(&TxNums, "tx", 100, "transaction nums")
	flag.IntVar(&SendTime, "st", 10, "send time")
	flag.StringVar(&CCkey, "cck", "default key", "chaincode key")
	flag.StringVar(&CCVal, "ccv", "default value", "chaincode value")
	flag.StringVar(&Opt, "opt", "invoke", "operation option")

	flag.Parse()

	// connect to fabric network
	connect.ConnectFabric(ChannelName, ChaincodeName)
	// init routepool
	routepool.Rpool.Init(MAX_ROUTINE_NUM)
	defer routepool.Rpool.Close()

	// init ledger before invoke transaction
	initLedger()

	// test: submit some transactions
	//submit.SimpleSubmitTxWithSR(100, "changeCarOwner", SendRate)
	//submit.InvokeTx(TxNums, "changeCarOwner")
	//submit.SubmitZipfianTxWithSR(TxNums, "changeCarOwner", SendRate)
	if ConflictRate != "init file" {
		submit.SubmitZipfianTxWithSRCR(SendRate, ConflictRate, SendTime)
	}
	if Opt == "create" && CCkey != "default key" {
		submit.CreateSingleSmallbankTx(SendRate, CCkey)
	}
	if Opt == "invoke" && CCkey != "default key" && CCVal != "default value" {
		submit.SubmitSingleSmallbankTx(SendRate, CCkey, CCVal)
	}

	fmt.Println("closing...")
}
