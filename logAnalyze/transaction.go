package logAnalyze

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: ygzhang
 * @Date: 2024/2/22 16:37
 * @Func: 根据日志计算出最终结果
 **/

type TransactionLatency struct {
	EndorsedTime string //背书时延
	OrderedTime  string //排序时延
	ValidateTime string //验证时延
}

type TransactionTime struct {
	CreateTime   time.Time //创建的交易的时间
	EndorsedTime time.Time //完成背书的时间
	OrderedTime  time.Time //完成排序的时间
	FinishTime   time.Time //交易结束的时间
}
type EndorseContext struct {
	TxId    string
	Time    string
	Latency string
}
type ReceiveContext struct {
	BlockId string
	Time    string
}
type ValidateContext struct {
	BlockId string
	Time    string
}
type CommitContext struct {
	BlockId string
	TxNums  int
	Time    string
}

// var EndorseFile = "./data/endorse.csv"
// var OrderedFile = "./data/receive"
// var ValidateFile = "./data/validate.csv"
// var CommitFile = "./data/commit.csv"
var EndorseFile = "data/endorse.csv"
var OrderedFile = "data/receive.csv"
var ValidateFile = "data/validate.csv"
var CommitFile = "data/commit.csv"

var EndorseText []EndorseContext //txid + endorse time
var ReceiveText []ReceiveContext
var ValidateText []ValidateContext
var CommitText []CommitContext

var TxTs []TransactionTime
var TxLs []TransactionLatency

const layout = "2006-01-02 15:04:05.000 MST" // Go的布局字符串

// ===========打印结果================= //
func WriteResult() (string, string, string) {
	var endorsedTimes, orderedTimes, validateTimes []string
	for _, latency := range TxLs {
		endorsedTimes = append(endorsedTimes, latency.EndorsedTime)
		orderedTimes = append(orderedTimes, latency.OrderedTime)
		validateTimes = append(validateTimes, latency.ValidateTime)
	}

	averageEndorsed := calculateAverage(endorsedTimes)
	averageOrdered := calculateAverage(orderedTimes)
	averageValidate := calculateAverage(validateTimes)

	return averageEndorsed, averageOrdered, averageValidate
}

// calculateAverage 计算字符串切片表示的数字的平均值，并返回平均值的字符串表示
func calculateAverage(strings []string) string {
	var total float64 = 0
	for _, str := range strings {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			// 错误处理：如果无法解析字符串为数字，可以打印错误并跳过该值
			fmt.Println("Error parsing string:", err)
			continue
		}
		total += num
	}
	average := total / float64(len(strings))
	return fmt.Sprintf("%.2f", average) // 返回保留两位小数的平均值
}

// ===========分析数据计算时间=============== //

func CalcLatency() {
	// 先计算Create, Endorse, Received Time
	CalcEVCTimeFromBlock()
	CalcCETime()
	// 计算各个阶段的latency
	for i := range TxTs {
		ct, et, ot, ft := TxTs[i].CreateTime, TxTs[i].EndorsedTime, TxTs[i].OrderedTime, TxTs[i].FinishTime
		txl := TransactionLatency{
			EndorsedTime: fmt.Sprintf("%d", et.Sub(ct).Milliseconds()),
			OrderedTime:  fmt.Sprintf("%d", ot.Sub(et).Milliseconds()),
			ValidateTime: fmt.Sprintf("%d", ft.Sub(ot).Milliseconds()),
		}
		TxLs = append(TxLs, txl)
	}
}

func CalcCETime() {
	txNums := len(TxTs)
	// 从文件中获取createTime, FinishedTime
	txs := loadTxsFromFile(txNums)
	for i := 0; i < txNums; i++ {
		TxTs[i].CreateTime = txs[i].CreateTime
		TxTs[i].FinishTime = txs[i].FinishTime
	}
}

// 根据block计算Endorse, Validate和Commit时间
func CalcEVCTimeFromBlock() {
	txns := 0
	for i := range CommitText {
		txn := CommitText[i].TxNums
		blockId := CommitText[i].BlockId
		txns += txn
		if ReceiveText[i].BlockId != blockId {
			log.Fatalln("blockId mistach receive text!")
		}
		if ValidateText[i].BlockId != blockId {
			log.Fatalln("blockId mistach validate text!")
		}
		for j := 0; j < txn; j++ {
			// 将字符串转time.Time类型
			rt, _ := time.Parse(layout, ReceiveText[i].Time)
			txt := TransactionTime{
				OrderedTime: rt,
			}
			TxTs = append(TxTs, txt)
		}
	}
	// 填充endorse time
	tlen := len(EndorseText)
	for i := tlen - 1; i >= tlen-txns; i-- {
		et, _ := time.Parse(layout, EndorseText[i].Time)
		TxTs[i-(tlen-txns)].EndorsedTime = et
	}
}

// =============读取文件到内存中=================//
func readEndorseFile() {
	// 打开文件
	file, err := os.Open(EndorseFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个扫描器来逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 读取当前行
		line := scanner.Text()
		// 按逗号分割行
		parts := strings.Split(line, ",")
		// 存储从每一行提取的数据
		if parts[1] == "0" { //特殊处理下为0的情况
			parts[1] = "1"
		}
		text := EndorseContext{
			Time:    parts[0],
			Latency: parts[1],
			TxId:    parts[2],
		}

		// 检查分割后的长度是否正确
		if len(parts) == 3 {
			EndorseText = append(EndorseText, text)
		}
	}
}
func readReceiveFile() {
	// 打开文件
	file, err := os.Open(OrderedFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个扫描器来逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 读取当前行
		line := scanner.Text()
		// 按逗号分割行
		parts := strings.Split(line, ",")
		// 存储从每一行提取的数据
		text := ReceiveContext{
			BlockId: parts[0],
			Time:    parts[1],
		}

		// 检查分割后的长度是否正确
		if len(parts) == 2 {
			ReceiveText = append(ReceiveText, text)
		}
	}
}
func readValidateFile() {
	// 打开文件
	file, err := os.Open(ValidateFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个扫描器来逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 读取当前行
		line := scanner.Text()
		// 按逗号分割行
		parts := strings.Split(line, ",")
		// 存储从每一行提取的数据
		text := ValidateContext{
			BlockId: parts[0],
			Time:    parts[1],
		}

		// 检查分割后的长度是否正确
		if len(parts) == 2 {
			ValidateText = append(ValidateText, text)
		}
	}
}
func readCommitFile() {
	// 打开文件
	file, err := os.Open(CommitFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个扫描器来逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 读取当前行
		line := scanner.Text()
		// 按逗号分割行
		parts := strings.Split(line, ",")
		// 存储从每一行提取的数据
		txNums, _ := strconv.Atoi(parts[1]) // 转换字符串为int，并处理可能的错误
		text := CommitContext{
			BlockId: parts[0],
			TxNums:  txNums,
			Time:    parts[2],
		}

		// 检查分割后的长度是否正确
		if len(parts) == 3 {
			CommitText = append(CommitText, text)
		}
	}
}
