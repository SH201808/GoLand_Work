package main

import (
	"../utils"
	"fmt"
	"os"
	"strconv"
	"time"
)

func IsExist()bool {
	_, err := os.Stat("E:\\Users\\27761\\GolandProjects\\awesomeProject\\work4\\mq\\mq1.mq")
	if err != nil {
		if os.IsExist(err){
			return true
		}else{
			return false
		}
	}
	return true
}

func oddDecode(data []byte)(ansData []byte,err error) {
	for i := 0; i < len(data); i++ {
		count := 0
		for j := 0; j < 8; j++ {
			temp := (data[i] >> j) & 1
			count += int(temp)	//	检验1的个数
		}
		if count %2 ==1 {
			ans :=data[i]>>1
			ansData = append(ansData,ans)
		}else{
			fmt.Println("发送错误")
			return
		}
	}
	return ansData,nil
}

func setExpression(data []byte,dataHash map[string]int)(name string) {
	flag := 0
	var s, Operator string
	var value int
	var tempData []byte
	//例：接受到的字符串为id_1 = 40
	for i := 0; i < len(data); i++ {
		if data[i] == 32 {	//当为空格时
			if flag == 0 {	//第一次遇到空格时，则tempData的数据为变量名
				s = string(tempData)
				tempData = tempData[0:0]	//清空tempData
			} else if flag == 1 {	//第二次遇到空格时，tempData的数据为运算符
				Operator = string(tempData)
				tempData = tempData[0:0]	//清空tempData
			}
			flag += 1
		} else {	//没有遇到空格
			if flag == 2 {	//遇到两次空格，tempData储存值
				for ; i < len(data); i++ {
					tempData = append(tempData, data[i])
				}
				//匹配运算符
				if Operator == "=" {
					value, _ = strconv.Atoi(string(tempData))
					dataHash[s] = value
					return s
				} else {
					value, _ = strconv.Atoi(string(tempData))
					switch Operator {
					case "+=":
						dataHash[s] += value
					case "-=":
						dataHash[s] -= value
					case "*=":
						dataHash[s] *= value
					case "/=":
						dataHash[s] /= value
					case "%=":
						dataHash[s] %= value
					}
				}
			} else {
				tempData = append(tempData, data[i])
			}
		}
	}
	return
}
func main() {
	finalData := make([]byte,0)
	tempData := make([]byte,0)
	nameData := make([]string,0)
	//dataHash := make(map[string]*expression)
	dataHash := make(map[string]int)	//储存表达式的变量和值
	ticker := time.NewTicker(1 * time.Second)	//一个定时器每一秒读取文件内容
	defer ticker.Stop()
	for range ticker.C {
		//如果文件不存在跳过
		if !IsExist() {
			continue
		}
		ReceiveData, err := utils.Decode()	//读取文件
		if err != nil {
			fmt.Println("Decode err = ", err)
			return
		}
		finalData, err = oddDecode(ReceiveData)	//对接受到的内容进行奇校验解码
		if err != nil {
			fmt.Println("oddDecode err = ", err)
			return
		}
		for i := 0; i < len(finalData); i++ {
			if finalData[i] == 92 {	//接受到的数据中有"" \n 等无效数据,需要跳过这些数据
				i += 4	//跳过无效数据
				name := setExpression(tempData, dataHash)	//将表达式的变量，值放入到map中，得到name
				nameData = append(nameData,name)	//将变量按输入顺序储存到切片中
				tempData = tempData[0:0]	//清空tempData
			} else {
				tempData = append(tempData, finalData[i])
			}
		}
		//按序输出表达式
		for i:= 0;i<len(dataHash);i++ {
			fmt.Printf(nameData[i]+" = "+"%d\n",dataHash[nameData[i]])
		}
		//读取完成后删除文件
		err = os.Remove("E:\\Users\\27761\\GolandProjects\\awesomeProject\\work4\\mq\\mq1.mq")
		if err != nil {
			fmt.Println("删除错误 err=", err)
			return
		}
	}
}