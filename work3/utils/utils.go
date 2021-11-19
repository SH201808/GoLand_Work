package utils

import (
	"encoding/binary"
	"fmt"
	"net"
)

func Encode(Conn net.Conn,UserData []byte)(err error){
	DataLen := uint32(len(UserData))	//得到Data长度
	var buf [1024]byte
	binary.BigEndian.PutUint32(buf[:4],DataLen)	//将长度放入buf中
	_,errWLen := Conn.Write(buf[:4])	//发送长度
	if errWLen != nil{
		fmt.Println("Conn.Write DataLen err= ",errWLen)
		return
	}
	n,errWData := Conn.Write(UserData)	//发送数据
	if n!= int(DataLen) || errWData != nil{
		fmt.Println("Conn.Write UserData err= ",errWData)
		return
	}
	return
}

func Decode(Conn net.Conn)(UserData []byte,err error){
	var buf [1024]byte
	n,errRLen := Conn.Read(buf[:4])	//读取发送的长度
	if n!=4||errRLen != nil{
		fmt.Println("Conn.ReadDataLen err= ",err)
		return
	}
	UserDataLen:=binary.BigEndian.Uint32(buf[:4])	//将发送的长度转为uint32
	_,errRData:=Conn.Read(buf[:UserDataLen])	//读取数据的byte切片
	if errRData != nil{
		fmt.Println("Conn.ReadData err= ",err)
		return
	}
	return buf[:UserDataLen],nil	//返回数据的切片
}