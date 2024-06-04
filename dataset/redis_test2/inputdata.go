// input data to redis-db

package main

import(
    "fmt"
	"github.com/go-redis/redis"
	
	"io"
	"bufio"
	"os"
	"strings"
	//"path/filepath"
	//"math/rand"

	"strconv"
)

var client *redis.Client

func initRedis()(err error){
	client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",  // 指定
		Password: "",
		DB:1,		// redis一共16个库，指定其中一个库即可
	})
    _,err = client.Ping().Result()
	return
}

func main() {
	err := initRedis()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n",err)
		return
	}
	
	path := "/home/Concurrency/fabric1_4/fabric-samples-1_4_2/fabric-samples/fabewallet/test/redis_test2/DataSet/"
	zipf := os.Args[1]
	rate := os.Args[2]
	fileNum := os.Args[3]

	filepath := path + "data" + zipf + "_" + rate + "/data" + fileNum

	fileHanle,err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
	}
 
	defer fileHanle.Close()
 
	reader := bufio.NewReader(fileHanle)

	// 按行处理
	for  {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		arr := strings.Split(string(line), " ")
		tx := ""
		if arr[0] == "0"{
			payer, remittee := arr[1], arr[2]
			tx = arr[0] + " owner" + payer + " owner" + remittee
		}else{
			tx = arr[0] + " owner" + arr[1]
		}
				
		client.LPush("transaction_queue", tx)
	}

	le, err := client.LLen("transaction_queue").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("current length: " + strconv.FormatInt(le, 10))
}

