package logAnalyze

import (
	"bufio"
	"io"
	"os"
	"regexp"
)

/**
 * @Author: ygzhang
 * @Date: 2024/2/22 16:34
 * @Func: 分析Docker 日志
 **/

// docker 导出日志的方法：docker logs --since 120m peer0.org1.example.com  >>input.log //将最近120min的日志导出到input.log中
// go run main.go 运行日志分析生成结果文件

var LogFile = "./dockerLog/logfile.log"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func writeValidateTime(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	var filename = "./data/validate.csv"
	var f_out *os.File
	var err1 error

	if checkFileIsExist(filename) { //如果文件存在
		f_out, err1 = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) //打开文件
		//fmt.Println("文件存在")
	} else {
		f_out, err1 = os.Create(filename) //创建文件
		//fmt.Println("文件不存在")
	}
	check(err1)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		//validate time
		vPattern := `Validated block.*\[([0-9]+)\] in ([0-9]+)ms\n$`
		result2 := substract(line, vPattern)
		if len(result2) > 0 {
			blockId := result2[1]
			validationTime := result2[2]

			wireteString := blockId + "," + validationTime + "\n"

			_, err1 := io.WriteString(f_out, wireteString) //写入文件(字符串)
			check(err1)

			//fmt.Println("blockId", blockId, "validationTime", validationTime)
		}
	}
}

func writeReceiveTime(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	var filename = "./data/receive.csv"
	var f_out *os.File
	var err1 error

	if checkFileIsExist(filename) { //如果文件存在
		f_out, err1 = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) //打开文件
		//fmt.Println("文件存在")
	} else {
		f_out, err1 = os.Create(filename) //创建文件
		//fmt.Println("文件不存在")
	}
	check(err1)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		//receive time
		//rPattern := `^.*([0-9][0-9]:[0-9][0-9]:[0-9][0-9]\.[0-9][0-9][0-9]).*Received block.*\[([0-9]+)\].*from buffer channel=mychannel\n$`
		rPattern := `^.*(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} UTC).*Received block.*\[([0-9]+)\].*from buffer channel=mychannel\n$`
		result1 := substract(line, rPattern)
		if len(result1) > 0 {
			receiveTime := result1[1]
			blockId := result1[2]

			writeString := blockId + "," + receiveTime + "\n"

			_, err1 := io.WriteString(f_out, writeString) //写入文件(字符串)
			check(err1)

			//fmt.Println("receiveTime", receiveTime, "blockId", blockId)
		}
	}
}

func writeCommitTime(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	var filename = "./data/commit.csv"
	var f_out *os.File
	var err1 error

	if checkFileIsExist(filename) { //如果文件存在
		f_out, err1 = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) //打开文件
		//fmt.Println("文件存在")
	} else {
		f_out, err1 = os.Create(filename) //创建文件
		//fmt.Println("文件不存在")
	}
	check(err1)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		//commit time
		//cPattern := `^.*(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} UTC).*Committed block.*\[([0-9]+)\].*with ([0-9]+) transaction\(s\).*in ([0-9]+)ms.*\(state_validation.*`
		cPattern := `Committed block.*\[([0-9]+)\].*with ([0-9]+) transaction\(s\).*in ([0-9]+)ms.*\(state_validation.*`

		result3 := substract(line, cPattern)
		if len(result3) == 4 {
			//commitTime := result3[1]
			blockId := result3[1]
			txNum := result3[2]
			commitTime := result3[3]

			wireteString := blockId + "," + txNum + "," + commitTime + "\n"

			_, err1 := io.WriteString(f_out, wireteString) //写入文件(字符串)
			check(err1)

			// fmt.Println("blockId",blockId,"txNum",txNum,"commitTime",commitTime)
		}
	}
}

func writeEndorseTime(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	var filename = "./data/endorse.csv"
	var f_out *os.File
	var err1 error

	if checkFileIsExist(filename) { //如果文件存在
		f_out, err1 = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) //打开文件
		//fmt.Println("文件存在")
	} else {
		f_out, err1 = os.Create(filename) //创建文件
		//fmt.Println("文件不存在")
	}
	check(err1)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		//endorse time
		// rPattern := `^.*([0-9][0-9]:[0-9][0-9]:[0-9][0-9]\.[0-9][0-9][0-9]).*Received block.*\[([0-9]+)\].*from buffer channel=mychannel\n$`
		rPattern := `^.*(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} UTC).*finished chaincode: \w+ duration: (\d)+ms.*channel=mychannel txID=([a-f0-9]+)`
		result1 := substract(line, rPattern)
		if len(result1) > 0 {
			endorseTime := result1[1]
			endorseLatency := result1[2]
			txID := result1[3]
			wireteString := endorseTime + "," + endorseLatency + "," + txID + "\n"

			_, err1 := io.WriteString(f_out, wireteString) //写入文件(字符串)
			check(err1)
		}
	}
}

func substract(content string, pattern string) []string {
	regR := regexp.MustCompile(pattern)
	if regR == nil {
		result := []string{}
		//fmt.Println("regexp err")
		return result
	}
	return regR.FindStringSubmatch(content)
}
