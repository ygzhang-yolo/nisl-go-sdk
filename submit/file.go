package submit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nisl-go-sdk/transactions"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/15 16:09
 * @Func:
 **/

//
// extractNumber
//  @Description: 提取出文件名的序号数字, 比如a.11.db, 提取出数字11
//  @param fileName
//  @return int
//
func extractNumber(fileName string) int {
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return 0
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return num
}

//
// getNextFileName
//  @Description: 返回要生成的最新一个next的名字
//  @param currentFile
//  @return string
//
func getNextFileName(currentFile string) string {
	sequenceNumber := extractNumber(currentFile) + 1
	return fmt.Sprintf("txs.%d.db", sequenceNumber)
}

//
// writeToFile
//  @Description: 将交易txs写入文件
//  @param dir
//  @param txs
//  @return error
//
func writeToFile(dir string, txs []transactions.Transaction) error {
	// 读取路径下的文件
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// 排序文件并找到最新的文件
	var dbFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".db") {
			dbFiles = append(dbFiles, file.Name())
		}
	}
	sort.Slice(dbFiles, func(i, j int) bool {
		return extractNumber(dbFiles[i]) < extractNumber(dbFiles[j])
	})
	var latestFile string
	if len(dbFiles) > 0 {
		latestFile = dbFiles[len(dbFiles)-1]
	} else {
		latestFile = "txs.1.db" // Start with txs.1.db
	}
	// 打开最新的文件
	filePath := filepath.Join(dir, latestFile)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	// 将一批交易写入文件txs
	for _, tx := range txs {
		// 将交易tx进行JSON序列化
		jsonData, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		jsonData = append(jsonData, '\n') // Append newline for each transaction

		// 检查文件大小是否超过100MB，超过需要创建新文件
		if fileInfo, _ := file.Stat(); fileInfo.Size()+int64(len(jsonData)) > 100*1024*1024 {
			file.Close() // Close the current file
			newFilePath := filepath.Join(dir, getNextFileName(latestFile))
			file, err = os.Create(newFilePath)
			if err != nil {
				return err
			}
			latestFile = newFilePath
		}

		// 将JSON数据写入文件
		if _, err = file.Write(jsonData); err != nil {
			file.Close()
			return err
		}
	}
	return nil
}

//
// writeToTempFile
//  @Description: 每次交易写入一个临时文件
//  @param dir
//  @param txs
//  @return error
//
func writeToTempFile(dir string, txs []transactions.Transaction) error {
	tempFile := "tx_temp.db"
	// 打开最新的文件
	filePath := filepath.Join(dir, tempFile)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	// 将一批交易写入文件txs
	for _, tx := range txs {
		// 将交易tx进行JSON序列化
		jsonData, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		jsonData = append(jsonData, '\n') // Append newline for each transaction
		// 将JSON数据写入文件
		if _, err = file.Write(jsonData); err != nil {
			file.Close()
			return err
		}
	}
	return nil
}
