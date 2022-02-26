package utils

import (
	"chatroomRedis/db"
	"chatroomRedis/model"
	"chatroomRedis/mq"
	"encoding/json"
	"github.com/gorilla/websocket"
)

func UserInRoom (roomId string,userId string)(flag bool,clients []model.Client,index,userNumber int,err error){
	reply,err := mq.Rdb.Get(roomId).Bytes()			//获取房间内的用户
	if err != nil {
		return false,nil,0,0,err
	}
	err = json.Unmarshal(reply,&clients)
	if err != nil {
		return false,nil,0,0,err
	}
	//判断用户是否在房间内
	for k,v := range clients{
		if v.UserId ==userId{
			return true,clients,k,len(clients),nil
		}
	}
	return false,clients,0,0,nil
}

func EnterRoom(userName,roomName string,wsconn *websocket.Conn)error{
	socketConn := model.WebsocketConn{
		UserName: userName,
		RoomName: roomName,
		Conn:   wsconn,
	}
	var socketConnes []*model.WebsocketConn
	i,ok:=model.ClientRooms.Load(userName)		//储存websocket连接
	if !ok{
		socketConnes = make([]*model.WebsocketConn,0)
	}else {
		socketConnes = i.([]*model.WebsocketConn)
	}
	socketConnes = append(socketConnes,&socketConn)
	model.ClientRooms.Store(userName, socketConnes)

	mq.Publish("system",roomName,userName+"成功进入房间")
	mq.CreateGroup(roomName,userName)			//创建消费者组绑定相应的队列
	//开启两个协程接受和读取消息
	go socketConn.Write()
	go socketConn.Read()
	return nil
}

func ExitRoom(clients []model.Client,index int,roomId,roomName,userId,userName string)error {
	clients = append(clients[:index], clients[index+1:]...)
	data, err := json.Marshal(clients)
	if err != nil {
		return err
	}
	err = mq.Rdb.Do("Set", roomId, data).Err()	//修改房间下的用户信息
	if err != nil {
		return err
	}
	user := model.User{}
	err = db.DB.Where("user_id = ?",userId).First(&user).Error		//判断用户是否在房间内
	if err != nil {
		return err
	}
	//修改用户下的房间信息
	reply,err := mq.Rdb.Get(user.UserName).Bytes()
	if err != nil {
		return err
	}
	var rooms []string
	json.Unmarshal(reply,&rooms)
	for k,v := range rooms{
		if v == roomId{
			rooms = append(rooms[:k],rooms[k+1:]...)
			break
		}
	}
	data,_ = json.Marshal(rooms)
	mq.Rdb.Do("Set",user.UserName,data)
	socketConnes, _ := model.ClientRooms.Load(userName)
	Connes := socketConnes.([]*model.WebsocketConn)
	for k, v := range Connes {
		if v.RoomName == roomName {
			v.Conn.Close()		//关闭websocket连接
			socketConnes = append(Connes[:k], Connes[k+1:]...)
			model.ClientRooms.Store(userName, socketConnes)

			mq.Publish("system",roomName,"exit"+userName)		//发送消息关闭协程
			return nil
		}
	}
	return nil
}
