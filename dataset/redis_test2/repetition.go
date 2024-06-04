package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

func calculateDuplicateRatio(lines []string) float64 {
	counts := make(map[string]int)
	totalWords := 0
	uniqueCount := 0

	for _, line := range lines {
		// 跳过空行
		if line == "" {
			continue
		}

		words := strings.Fields(line)
		totalWords += len(words)

		// 统计每个单词出现的次数
		for _, word := range words {
			counts[word]++
		}
	}

	// 统计没有重复出现的单词数量
	for _, count := range counts {
		if count == 1 {
			uniqueCount++
		}
	}

	duplicateRatio := 1.0 - (float64(uniqueCount) / float64(totalWords))
	return duplicateRatio
}

func main() {
	filename := os.Args[1] // 获取文件名参数
	batchSize := 100     // 每1000行单独计算重复率

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	numBatches := int(math.Ceil(float64(len(lines)) / float64(batchSize)))

	totalDuplicateRatio := 0.0

	for i := 0; i < numBatches; i++ {
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize

		if endIndex > len(lines) {
			endIndex = len(lines)
		}

		batchLines := lines[startIndex:endIndex]

		if len(batchLines) > 0 {
			duplicateRatio := calculateDuplicateRatio(batchLines)
			fmt.Printf("前 %d 行的重复率为：%.2f%%\n", endIndex, duplicateRatio*100)
			totalDuplicateRatio += duplicateRatio
		} else {
			fmt.Printf("前 %d 行的重复率为：NaN%%\n", endIndex)
		}
	}

	averageDuplicateRatio := totalDuplicateRatio / float64(numBatches)
	fmt.Printf("整个文件的重复率平均值为：%.2f%%\n", averageDuplicateRatio*100)
}
