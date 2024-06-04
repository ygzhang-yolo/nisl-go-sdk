package main

import (
	"bufio"
	"nisl-go-sdk/submit"
	"os"
	"strconv"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/12 14:33
 * @Func:
 **/

const ledgerPath = "data/ledger/ledger.data"

//
// initLedger
//  @Description: init ledger, create car
//
func initLedger() {
	// 检查ledger文件是否存在, 已经存在则不需要创建
	if _, err := os.Stat(ledgerPath); !os.IsNotExist(err) {
		return
	}
	//ledger不存在, 需要创建账本
	txNums := 10000
	file, _ := os.Create(ledgerPath)
	defer file.Close()
	// 写入文件
	writer := bufio.NewWriter(file)
	for i := 0; i < txNums; i++ {
		//writer.WriteString("CAR" + strconv.Itoa(i) + ", VW, Polo, Grey, Mary\n")
		writer.WriteString("owner" + strconv.Itoa(i) + ", 10000000\n")
	}
	//submit.CreateLedger("createCar", txNums)
	submit.CreateLedger("CreateAccount", txNums)
}

//result, err := contract.SubmitTransaction("CreateAccount", "owner"+strconv.Itoa(i) , "10000000")
