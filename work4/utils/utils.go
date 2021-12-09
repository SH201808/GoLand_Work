package utils

import (
	"encoding/binary"
	"fmt"
	"os"
)

func Encode (data []byte) error {
	var buf [1024]byte
	filePath := "./mq/mq1.mq"
	f, err := os.OpenFile(filePath,os.O_CREATE|os.O_WRONLY,0666)	//打开文件，文件不存在时可以创建文件并可写
	if err != nil {
		fmt.Println("open file err ")
		return err
	}
	defer f.Close()
	dataLen := len(data)	//数据长度
	binary.BigEndian.PutUint32(buf[:4],uint32(dataLen))
	_,err = f.Write(buf[:4])	//写入数据长度
	if err != nil {
		fmt.Println("发送长度错误",err)
		return err
	}
	_,err = f.Write(data)	//写入数据
	if err != nil {
		fmt.Println("发送数据错误",err)
		return err
	}
	return nil
}

func Decode() (data []byte,err error) {
	var buf [1024]byte
	f, err := os.Open("E:\\Users\\27761\\GolandProjects\\awesomeProject\\work4\\mq\\mq1.mq")	//打开文件
	if err != nil {
		fmt.Println("Openfile err = ", err)
		return
	}
	defer f.Close()
	if err != nil {
		fmt.Println("ReadAtLeast err = ", err)
		return
	}
	_,err = f.Read(buf[:4])	//读取文件的长度
	dataLen := binary.BigEndian.Uint32(buf[:4])
	data = make([]byte,dataLen)	//根据长度定义切片
	_,err = f.Read(data)	//读取文件
	if err != nil {
		fmt.Println("Read file err = ",err)
		return
	}
	return data,nil
}

