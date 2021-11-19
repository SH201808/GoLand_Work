package main

import (
	"awesomeProject/work3/User"
	"awesomeProject/work3/utils"
	"encoding/json"
	"fmt"
	"net"
)




func GetFind(Conn net.Conn){
	defer Conn.Close()	//延时关闭Conn
	Data := make(map[int64]string)	//服务器中存储的数据
	Data[123456] = "a"
	Data[12345] = "b"
	Data[1234] = "c"
	var UserDataId User.GetUserReq
	for {
		UserIdByte, errDecode := utils.Decode(Conn)	//通过Decode得到客户端传入的UseId的byte切片
		if errDecode != nil {
			fmt.Println("Decode err= ", errDecode)
			return
		}
		//如果从客户端传入的UserId的byte切片为nil,说明客户端未传入UserId的数据，说明客户端退出，则服务器端也退出
		if UserIdByte == nil{
			fmt.Println("客户端退出，服务器端也退出")
			return
		}
		errMarshal :=json.Unmarshal(UserIdByte, &UserDataId)	//把客户端传入的数据反序列化到结构体中
		if errMarshal !=nil{
			fmt.Println("Marshal err= ",errMarshal)
			return
		}
		UserId := UserDataId.UserID	//得UserId
		var UserResp User.GetUserResp
		_,ok:=Data[UserId]	//查找是否有UserId的值
		if ok {
			UserResp.UserID = UserId
			UserResp.UserName = Data[UserId]
			UserDataByte, errMarshal := json.Marshal(UserResp)	//将UserData的数据序列化
			if errMarshal != nil {
				fmt.Println("Marshal err= ", errMarshal)
				return
			}
			errEncode := utils.Encode(Conn, UserDataByte)	//将UserData数据发送到客户端
			if errEncode != nil {
				fmt.Println("Encode err= ", errEncode)
				return
			}
		}else{
			NotFind := "未找到"
			errEncode := utils.Encode(Conn,[]byte(NotFind))	//把未找到的切片发送到客户端
			if errEncode != nil {
				fmt.Println("Encode err= ", errEncode)
				return
			}
		}
	}
}

func main(){
	listen,err := net.Listen("tcp","localhost:8888")	//监听本地的8888端口
	if err != nil{
		fmt.Println("net.Listen err= ",err)
		return
	}
	defer listen.Close()  //关闭listen
	//循环等待客户端连接
	for{
		Conn,err := listen.Accept()
		if err != nil{
			fmt.Println("listen.Accept err= ",err)
			return
		}
		go GetFind(Conn)	//开启一个协程为客户端服务
	}
}