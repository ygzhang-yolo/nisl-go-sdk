package logAnalyze

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nisl-go-sdk/transactions"
	"os"
)

/**
 * @Author: ygzhang
 * @Date: 2024/2/20 12:41
 * @Func: 从tx1.db中读取交易的createTime
 **/

var filePath = "/home/zhangyiguang/fabric/nisl-go-sdk/data/temp/tx_temp.db"

func loadTxsFromFile(txNums int) []transactions.Transaction {
	// 加载倒数txNUms个数据到切片中
	file, err := os.Open(filePath) // 替换 "yourfile.json" 为你的文件名
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var transactions []transactions.Transaction

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var tx transactions.Transaction
		err := json.Unmarshal([]byte(scanner.Text()), &tx)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue // 如果有错误，跳过这行
		}
		transactions = append(transactions, tx)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	// 只保留最后txNums个JSON数据
	var lastTransactions []transactions.Transaction
	if len(transactions) > txNums {
		lastTransactions = transactions[len(transactions)-txNums:]
	} else {
		lastTransactions = transactions
	}
	return lastTransactions
}
