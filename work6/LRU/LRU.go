package LRU

import (
	"LRU/utils"
	"bufio"
	"container/list"
	"net"
	"strings"
	"sync"
)


type LRUCache struct {
	Cap int64				//容量
	LinkList *list.List		//双向链表
	Hash sync.Map			//哈希表
	Listener net.Listener
	Conn net.Conn
}

func(L *LRUCache)Run(){
	defer L.Listener.Close()
	for{
		L.Conn,_ = L.Listener.Accept()	//建立连接
		go process(L)					//处理连接
	}
}

func process(L *LRUCache){
	defer L.Conn.Close()
	for {
		msg, err := L.ReceiveMsg()		//接受信息
		if err != nil {
			L.SendMsg("接受数据错误")
			return
		}
		s := strings.Split(msg, " ")
		if s[0] == "Set" {
			L.update(s[1], s[2])		//更新缓存
			err = L.SendMsg("success")
		} else if s[0] == "Get" {
			msg = L.find(s[1])			//查询缓存
			if msg == "LRU中无记录" {
				err = L.SendMsg("LRU中无记录")
			} else {
				err = L.SendMsg(msg)
				L.update(s[1], msg)		//更新缓存
			}
		} else if s[0] == "Delete" {
			L.Hash.Delete(s[1])			//删除缓存
			L.SendMsg("success")
		}
	}
}

func (L *LRUCache)Get(key string)(string,error) {
	err := L.SendMsg("Get", key)		//发送信息
	if err != nil {
		return "发生信息错误", err
	}
	msg, err := L.ReceiveMsg()				//接受信息
	if err != nil {
		return "接受信息错误", err
	}
	return msg, nil
}

func (L *LRUCache)Set(key string, value string) error {
	err := L.SendMsg("Set", key, value)		//发送信息
	if err != nil {
		return err
	}
	msg, err := L.ReceiveMsg()					//接受信息
	if err != nil || msg != "success" {
		return err
	}
	return nil
}

func (L *LRUCache)Delete(key string) error{
	err := L.SendMsg("Delete",key)			//发送信息
	if err != nil {
		return err
	}
	_,err = L.ReceiveMsg()				//接受信息
	if err != nil {
		return err
	}
	return nil
}

func (L *LRUCache)update(key string,value string) {
	v,ok := L.Hash.Load(key)			//缓存储存新数据
	if ok  {
		L.LinkList.MoveToFront(v.(*list.Element))	//将缓存中已存在的数据移动到最左侧
		return
	}
	if L.LinkList.Len() == int(L.Cap) {
		last := L.LinkList.Back()
		L.LinkList.Remove(last)			//删除最后一个数据
	}
	e := L.LinkList.PushFront(value)	//将不存在的数据放到链表最左侧
	L.Hash.Store(key,e)					//储存新数据
}

func (L *LRUCache)find(key string)(msg string){
	v,ok := L.Hash.Load(key)		//查询缓存中是否有数据
	if !ok{
		msg = "LRU中无记录"
		return
	}
	msg = v.(*list.Element).Value.(string)
	return
}

func (L *LRUCache)SendMsg(x ...string)error{
	method := x[0]
	var s string
	//判断客户端发送的方法
	if method=="Set"{
		s = x[0]+" "+x[1]+" "+x[2]
	}else if method == "Get" || method=="Delete" {
		s = x[0] + " " + x[1]
	}else{
		s = x[0]
	}
	data,err := utils.Encode(s)		//将消息编码
	if err != nil {
		return err
	}
	L.Conn.Write(data)			//写入消息
	return nil
}

func(L *LRUCache)ReceiveMsg()(string,error){
	reader := bufio.NewReader(L.Conn)
	msg,err:=utils.Decode(reader)		//解码消息
	if err != nil {
		return "", err
	}
	return msg,nil
}