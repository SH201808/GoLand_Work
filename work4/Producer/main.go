package main

import (
	"../utils"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func oddCode(data uint8) (ans byte) {
	count := 0
	//判断每一个字节二进制的1的个数
	for i := 0; i <= 7; i++ {
		temp := (data >> i) & 1
		count += int(temp)
	}
	ans = data << 1 //留出校验位
	if count%2 == 0 {
		ans += 1
	} else {
		ans += 0
	}
	return ans
}


func main() {
	for {
		flag := 1
		fmt.Println("请输入计算式(按 id_1 = 40 这种格式输入)，或按0退出")
		var dataString []string //储存表达式字符串
		inputReader := bufio.NewReader(os.Stdin)
		for {
			input, err := inputReader.ReadString('\n') //换行时读取一行字符串
			if err != nil {
				fmt.Println("ReadString err = ", err)
				return
			}
			//退出
			if input == "0\n" {
				flag = 0
				break
			}
			dataString = append(dataString, input) //添加到切片中
		}
		data, err := json.Marshal(dataString) //对数据进行序列化
		var sendData []byte
		if err != nil {
			fmt.Println("err = ", err)
		}
		//对数据进行奇校验编码
		for i := 2; i < len(data)-2; i++ {
			ans := oddCode(data[i])
			sendData = append(sendData, ans) //结果添加到切片中
		}
		err = utils.Encode(sendData) //进行Encode协议写入到文件中
		if err != nil {
			fmt.Println("Encode err = ", err)
			return
		}
		if flag == 0 {
			break
		}
	}
}

