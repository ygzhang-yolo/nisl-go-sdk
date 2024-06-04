package transactions

import "time"

/**
 * @Author: ygzhang
 * @Date: 2024/1/11 14:04
 * @Func:
 **/

//定义一个交易
type Transaction struct {
	Id         int
	Status     TxStatus //交易状态，0-未处理, 1-成功, -1-失败
	Key        string
	Value      string
	CreateTime time.Time
	FinishTime time.Time
	CCFunc     string //交易对应的链码名
}

type TxStatus int

const (
	NOTSUBMIT TxStatus = iota
	SUCCESS
	FAIL
)

func (t Transaction) GetSendRate(txs []Transaction) float64 {
	transactionNums := len(txs)
	first_cr_time := txs[0].CreateTime
	last_cr_time := txs[0].CreateTime
	for _, tx := range txs {
		if tx.CreateTime.After(last_cr_time) {
			last_cr_time = tx.CreateTime
		}
		if tx.CreateTime.Before(first_cr_time) {
			first_cr_time = tx.CreateTime
		}
	}
	submitTime := last_cr_time.Sub(first_cr_time).Seconds()
	return float64(transactionNums) / submitTime
}

func (t Transaction) GetThroughput(txs []Transaction) float64 {
	transactionNums := len(txs)
	first_fin_time := txs[0].FinishTime
	last_fin_time := txs[0].FinishTime
	for _, tx := range txs {
		if tx.FinishTime.After(last_fin_time) {
			last_fin_time = tx.FinishTime
		}
		if tx.FinishTime.Before(first_fin_time) {
			first_fin_time = tx.FinishTime
		}
	}
	executeTime := last_fin_time.Sub(first_fin_time).Seconds()
	return float64(transactionNums) / executeTime
}

func (t Transaction) GetAvglatency(txs []Transaction) float64 {
	transactionNums := len(txs)
	var total_latency float64 = 0
	for _, tx := range txs {
		tx_latency := tx.FinishTime.Sub(tx.CreateTime).Seconds()
		total_latency += tx_latency
	}
	return total_latency / float64(transactionNums)
}

func testSSLSecurity() error {
	//
	//
	return nil
}