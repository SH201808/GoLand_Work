package main

import (
	"LRU/LRU"
	"container/list"
	"net"
	"strconv"
)
type Server interface {
	Run()
}

func NewServer(port int, maxSize int64) Server {
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(port)) //监听端口
	L := LRU.LRUCache{
		Cap:      maxSize,
		Listener: listener,
		LinkList: list.New(),
	}
	return &L
}

func main() {
	//一个只能装5条数据的 LRU 数据库, 服务端口 8080
	s := NewServer(8080, 2)
	//LRU 数据库运行
	s.Run()
}
