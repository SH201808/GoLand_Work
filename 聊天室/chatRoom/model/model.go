package model

import (
	"chatroomRedis/mq"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var ClientRooms sync.Map

type User struct {
	UserId    string	`gorm:"primary_key;AUTO_INCREMENT"form:"userId"`
	UserPwd   string	`form:"userPwd"`
	UserName  string	`form:"userName"`
	UserEmail string	`gorm:"unique"form:"userEmail"`
}

type Room struct {
	RoomId     string	`gorm:"primary_key"form:"roomId"`
	RoomName   string	`form:"roomName"`
	RoomCap    int		`form:"roomCap"`
	RoomAccess bool		`form:"roomAccess"`
	RoomOwner  string
	CreatedAt time.Time
}

type Client struct {
	UserId string
	UserName string
}

type WebsocketConn struct {
	UserName	string
	RoomName    string
	Conn *websocket.Conn
}

func (socketConn *WebsocketConn)Write(){
	for {
		_, message, err := socketConn.Conn.ReadMessage()		//从连接中得到发送的消息
		if err != nil {
			socketConn.Conn.Close()
			return
		}
		if len(message) == 0 {
			continue
		}
		mq.Publish(socketConn.UserName,socketConn.RoomName,string(message))		//将消息发送给群里所有成员
	}
}

func (socketConn *WebsocketConn)Read() {
	for {
		roomName, userName, content, id, err := mq.GetMsg(socketConn.UserName, socketConn.RoomName)		//从MQ中得到消息
		if err != nil {
			fmt.Println("得到消息", err)
			return
		}
		if content == "exit"+socketConn.UserName && userName == "system" {		//系统消息表示该用户退出房间
			flag := false
			connes, ok := ClientRooms.Load(socketConn.UserName)				//判断服务器中有无该用户的websocket连接
			if !ok{
				err = mq.Rdb.XGroupDestroy(socketConn.RoomName, socketConn.UserName).Err()		//删除MQ中的该消费者组
				mq.Rdb.XDel(socketConn.RoomName, id)
				if err != nil {
					fmt.Println("关闭协程", err)
					return
				}
				return
			}
			for _, v := range connes.([]*WebsocketConn) {
				if v.Conn == socketConn.Conn {
					flag = true
				}
			}
			//没有旧的websocket连接关闭协程
			if !flag {
				err = mq.Rdb.XGroupDestroy(socketConn.RoomName, socketConn.UserName).Err()
				mq.Rdb.XDel(socketConn.RoomName, id)
				if err != nil {
					fmt.Println("关闭协程", err)
					return
				}
				return
			}
		}else {
			msg := mq.Message{
				FromUserName: userName,
				RoomName:     roomName,
				Content:      content,
			}
			err = socketConn.Conn.WriteJSON(msg)		//向用户发送消息
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
