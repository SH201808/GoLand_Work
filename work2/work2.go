package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
)


var CpuDataList []float64	//储存CPU利用率数据
var MemUseDataList []int64	//储存内存使用量数据

var MaxCpu float64	//最大CPU利用率
var MinCpu  float64	//最小CPU利用率
var SumCpu float64	//CPU利用率总和
var MaxMemUse int64	//最大内存使用量
var MinMemUse int64	//最小内存使用量
var SumMemUse int64	//内存使用量总和


func init(){	//初始化变量
	MaxCpu = 0.0
	MinCpu = 100.0
	SumCpu = 0.0
	MaxMemUse = 0
	MinMemUse = int64(math.Pow(2,32)-1)
	SumMemUse = 0
}

func main(){
	data,err := os.Open("C:\\Users\\27761\\Desktop\\work2test10.txt")	//打开文件
	if err != nil{
		fmt.Println("open file err = ",err)	//如果文件打开错误，输出错误
	}
	defer data.Close()	//程序最后关闭文件
	CpuStrReg := regexp.MustCompile("\\d{0,3}.\\d.id")	//在文件中匹配CPU未利用率字符串 如：100.0 id
	CpuDataReg := regexp.MustCompile("\\d{0,3}.\\d")	//在CPU未利用率字符串中匹配有效数据 如：100.0
	MemStrlineReg := regexp.MustCompile("KiB Mem")		//在文件中匹配内存数据的那一行 如：KiB Mem :  3880172 total,  2839192 free,   246500 used,   794480 buff/cache
	MemUseDataStrReg := regexp.MustCompile("[0-9]+.used")	//在Mem行中匹配内存使用量数据字符串 如：246500 used
	MemUseDataReg := regexp.MustCompile("[0-9]+")	//匹配内存有效数据 如：246500
	if CpuStrReg == nil || CpuDataReg == nil || MemStrlineReg ==nil || MemUseDataStrReg == nil||MemUseDataReg == nil{
		fmt.Println("MustComplie err")	//如果compile错误，输出错误
	}
	reader :=bufio.NewReader(data) //创建缓冲
	//读取文件内容
	for{
		Str,err := reader.ReadString('\n')	//读取文件的每一行
		if err == io.EOF {
			break	//读到文件末尾退出
		}
		CpuStr := CpuStrReg.FindString(Str)	//匹配CPU未利用率字符串
		MemStrline := MemStrlineReg.FindString(Str)	//匹配内存数据的那一行
		if CpuStr != ""{
			CpuData,_ := strconv.ParseFloat(CpuDataReg.FindString(CpuStr),64)	//得到CPU未利用率
			CpuData = 100.0-CpuData		//得到CPU利用率
			//更新最大值，最小值，总量
			if MaxCpu < CpuData{
				MaxCpu = CpuData
			}else if MinCpu > CpuData{
				MinCpu = CpuData
			}
			SumCpu += CpuData
			CpuDataList = append(CpuDataList,CpuData)	//将CPU利用率添加到切片中
		}else if MemStrline != ""{
			MemUseDataStr := MemUseDataStrReg.FindString(Str)	//匹配内存使用量数据字符串
			MemUseData,_ := strconv.ParseInt(MemUseDataReg.FindString(MemUseDataStr),10,64)	//匹配内存使用量数据
			MemUseData /= 1000	//设单位为MB
			////更新最大值，最小值，总量
			if MaxMemUse < MemUseData{
				MaxMemUse = MemUseData
			}else if MinMemUse > MemUseData{
				MinMemUse = MemUseData
			}
			SumMemUse += MemUseData
			MemUseDataList = append(MemUseDataList,MemUseData)	//将内存使用量添加到切片中
		}
	}
	AvgCpu := SumCpu/float64(len(CpuDataList))	//计算CPU平均利用率
	fmt.Println("Cpu利用率")
	for _,v := range CpuDataList{
		fmt.Printf("%.1f\n",v)	//输出CPU利用率
	}
	fmt.Printf("Max:%.1f Min:%.1f ave:%.1f",MaxCpu,MinCpu,AvgCpu)	//输出最大值，最小值，平均值

	AvgMemUse := SumMemUse/int64(len(MemUseDataList))	//计算内存平均使用量
	fmt.Println("\n内存使用量")
	for _,v := range MemUseDataList{
		fmt.Printf("%d\n",v)	//输出内存使用量
	}
	fmt.Printf("MaxMemUse:%d MinMemUse:%d AvgMemUse:%d",MaxMemUse,MinMemUse,AvgMemUse)	//输出最大值，最小值，平均值
}