package main

import (
	"awesomeProject/work3/User"
	"awesomeProject/work3/utils"
	"encoding/json"
	"fmt"
	"net"
)



func main(){
	Conn,errDial:=net.Dial("tcp","localhost:8888")	//连接本地的8888端口
	if errDial != nil{
		fmt.Println("net.Dial err=",errDial)
		return
	}
	defer Conn.Close()	//程序结束后关闭Conn
	var key int
	//循环进行用户查找
	for{
		fmt.Println("1.查找用户名")
		fmt.Println("2.退出系统")
		fmt.Println("请按序号选择操作")
		fmt.Scanf("%d\n",&key)
		if key == 1{
			var userId int64
			fmt.Println("请输入用户ID")
			fmt.Scanf("%d\n",&userId)
			User := &User.GetUserReq{
				UserID: userId,
			}
			UserData,errMarshal:=json.Marshal(User)	//对User进行序列化得到字节数组
			if errMarshal != nil{
				fmt.Println("json.Marshal err=",errMarshal)
				return
			}
			errEncode := utils.Encode(Conn,UserData)	//通过utils中的Encode向服务器传输数据
			if errEncode != nil {
				fmt.Println("utils.Encode err=", errEncode)
				return
			}
			UserDataRes, errDecode := utils.Decode(Conn)	//通过utils中的Decode得到服务器传回的数据
			if errDecode != nil {
				fmt.Println("utils.Decode err=", errDecode)
				return
			}
			fmt.Println("Res:", string(UserDataRes))
		}else if key == 2{
			fmt.Println("退出系统")
			return
		}else{
			fmt.Println("输入错误请重新输入")
		}
	}
}