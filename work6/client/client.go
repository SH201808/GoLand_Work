package client

import (
	"LRU/LRU"
	"net"
)

type Client interface {
	//Set() 如果没有就新建,有就更新, 失败返回错误
	Set(key string, value string) error
	//Get() 获取缓存的内容, 失败返回错误
	Get(key string) (string, error)
	//Delete() 删除缓存的内容, 失败返回错误
	Delete(key string) error
}

func NewClient(addr string) Client {
	conn, _ := net.Dial("tcp", addr)
	L := LRU.LRUCache{
		Conn: conn,
	}
	return &L
}

