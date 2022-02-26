package mq

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Message struct {
	FromUserName  string
	RoomName    string
	Content     string
}

func Publish(userName,roomName,content string){
	a := redis.XAddArgs{
		Values: make(map[string]interface{},0),		//用户消息
		Stream: roomName,							//发送到的房间
		ID: "*",
	}
	a.Values[userName] = content					//消息内容
	err := Rdb.XAdd(&a).Err()						//将消息发送到MQ
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateGroup(roomName,userName string)error{
	err := Rdb.XGroupCreate(roomName,userName,"0").Err()		//创建消费者组，将消费者绑定到相应的队列
	if err != nil {
		return err
	}
	return nil
}

func GetMsg(userId,roomId string)(roomName,userName,content,id string,err error){
	var stream []string
	stream = append(stream,roomId,">")				//储存房间号
	a := redis.XReadGroupArgs{
		Group: userId,								//消费者组
		Consumer: userId,							//消费者
		Streams: stream,							//房间号
		Count: 1,
		Block: 0,									//阻塞，一直到有消息
	}
	x,err := Rdb.XReadGroup(&a).Result()			//读取消息
	if err != nil {
		return "", "", "", "", err
	}
	for _,msgs := range x{
		roomName = msgs.Stream
		for _,value := range msgs.Messages{
			for k,msgstring := range value.Values{
				userName = k
				return roomName,userName,msgstring.(string),value.ID, err
			}
		}
	}
	return "", "", "", "", nil
}