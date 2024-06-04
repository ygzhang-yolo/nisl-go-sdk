package main

import (
	// "bufio"
	"fmt"
	"os"
	"sort"
	"io/ioutil"
	"strconv"
	"strings"
	"math"
	"math/rand"
	"time"
)

func main() {

	filename := os.Args[1] // 获取文件名参数
	rateStr := os.Args[2] // 获取读写集比例，例如输入10，则代表读集10%
	outputfilename := os.Args[3] // 获取输出文件名参数

	batchSize := 500     // 每1000行单独计算重复率

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	numBatches := int(math.Ceil(float64(len(lines)) / float64(batchSize)))

	rate, err := strconv.Atoi(rateStr)
	if err != nil {
		fmt.Println("Error converting tag to integer:", err)
		return
	}
	rate = rate * 5

	for i := 0; i < numBatches; i++ {
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize

		if endIndex > len(lines) {
			endIndex = len(lines)
		}

		batchLines := lines[startIndex:endIndex]

		arrayRand := [500]int{0}
		updateArray(&arrayRand, rate)
		
		counts := make(map[string]int)
		calculateDuplicate(&counts, batchLines)

		settleLines(&counts, &arrayRand, batchLines, outputfilename)
	}
}

func settleLines(counts *map[string]int, arrayRand *[500]int, lines []string, outputfilename string) {
	// 打开或创建输出文件，使用追加模式
	outputFile, err := os.OpenFile(outputfilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening or creating output file:", err)
		return
	}
	defer outputFile.Close()

	tag := -1

	for _, line := range lines {
		tag++
		
		// 跳过空行
		if line == "" {
			continue
		}

		result := line

		for i := 1; i <= 10000; i++{
			account := strconv.Itoa(i)

			if (*counts)[account] > 0{
				continue
			}else{
				(*counts)[account]++
				result = fmt.Sprintf("%s %s",line,account)
				break
			}
		}

		// 在这里写入文件，例如：
		outputLine := fmt.Sprintf("%s\n",result)
		_, err = outputFile.WriteString(outputLine)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

// 选择出 count 个不同的在 0 到 max 之间的整数
func selectRandomNumbers(max, count int) []int {
	if count > max {
		count = max
	}

	// 初始化可选的整数列表
	availableNumbers := make([]int, max)
	for i := 0; i < max; i++ {
		availableNumbers[i] = i
	}

	// 随机选择 count 个不同的整数
	selectedNumbers := make([]int, count)
	for i := 0; i < count; i++ {
		// 使用当前时间的Unix时间戳作为种子
		rand.Seed(time.Now().UnixNano())

		// 从可选的整数中随机选择一个
		index := rand.Intn(max - i)
		selectedNumbers[i] = availableNumbers[index]
		// 将已选的整数从可选列表中移除
		availableNumbers[index], availableNumbers[max-i-1] = availableNumbers[max-i-1], availableNumbers[index]
	}

	return selectedNumbers
}

// 随机选择一个在 min 到 max 之间的正整数
func getRandomNumber(min, max int) int {
	if min > max {
		return 0
	}

	// 使用当前时间的Unix时间戳作为种子
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min+1) + min
}

// 从给定的切片中随机挑选 n 个元素
func getRandomSubset(numbers []int, n int) []int {
	if n > len(numbers) {
		return nil
	}

	// 使用当前时间的Unix时间戳作为种子
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	return numbers[:n]
}

func updateArray(a *[500]int, rate int) {
	// 随机选择一个(rate,500)的正整数
	randomNumber := getRandomNumber(rate, 500)

	// 选择出 500 以内的 randomNumber 个不同的整数
	selectedNumbers := selectRandomNumbers(500, randomNumber)

	// 排序结果
	sort.Ints(selectedNumbers)

	// 从选定的数字中再随机挑选 rate 个
	selectedSubset := getRandomSubset(selectedNumbers, rate)

	for i := range *a {
		if contains(selectedNumbers, i+1) {
			(*a)[i] = 1
			if contains(selectedSubset, i+1) {
				(*a)[i] = 2
			}
		} else {
			(*a)[i] = 0
		}
	}
}

// 判断切片中是否包含某个元素
func contains(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func calculateDuplicate(counts *map[string]int, lines []string) {
	for _, line := range lines {
		// 跳过空行
		if line == "" {
			continue
		}

		words := strings.Fields(line)

		// 统计每个单词出现的次数
		for _, word := range words {
			(*counts)[word]++
		}
	}
}