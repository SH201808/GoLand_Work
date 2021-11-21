package main

import (
	"fmt"
	"runtime"
)

func adddata(datachan chan string,data map[int]string){
	for i:=0;i<40;i++{
		datachan <-data[i%4]
	}
	close(datachan)		//关闭datachan
}
func printname(datachan chan string,exitchan chan bool){
	for {
		s, ok := <-datachan
		if !ok {
			break
		}
		fmt.Println(s)
	}
	exitchan <- true	//表示该goroutine完成
}
func main(){
	runtime.GOMAXPROCS(1)	//只使用单核
	//将名字放入地图中
	data := make(map[int]string, 10)
	data[0] = "张三"
	data[1] = "李四"
	data[2] = "王五"
	data[3] = "赵六"
	datachan := make(chan string,40)	//存放数据的通道
	exitchan := make(chan bool, 3)		//缓冲通道
	go adddata(datachan,data)	//开启一个goroutine把名字依次放入chan中
	for i:=0;i<40;i++{
		datachan <-data[i%4]
	}
	close(datachan)
	//开启三个goroutine从datachan中接受名字并打印
	for j:=0 ; j<3 ;j++{
		go printname(datachan,exitchan)
	}
	for k:= 0 ;k<3 ;k++ {
		<-exitchan	//堵塞主线程，等待四个goroutine跑完
	}
}
