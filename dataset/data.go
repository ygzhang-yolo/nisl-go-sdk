package dataset

import (
	"bufio"
	"log"
	"os"
	"strings"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/12 16:25
 * @Func:
 **/

//
// GetDataFromFile
//  @Description: 从文件中读取键值, 以一个string切片返回
//  @param path
//  @return []string
//
func GetDataFromFile(path string) []string {
	var result []string
	// open file
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal("No such file,", path, err)
		return []string{}
	}
	// create a scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fileds := strings.Fields(line)
		// 只需要每行的第一个string
		result = append(result, fileds[0])
		//for _, str := range fileds {
		//	result = append(result, str)
		//}
	}
	return result
}

//
// GetAllDataFromFile
//  @Description: 从文件中读取键值, 以一个[][]string返回
//  @param path
//  @return []string
//
func GetAllDataFromFile(path string) [][]string {
	var result [][]string
	// open file
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal("No such file,", path, err)
		return [][]string{}
	}
	// create a scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fileds := strings.Fields(line)
		// 只需要每行的第一个string
		result = append(result, fileds)
		//for _, str := range fileds {
		//	result = append(result, str)
		//}
	}
	return result
}
