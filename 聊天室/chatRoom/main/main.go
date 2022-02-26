package main

import (
	"chatroomRedis/db"
	"chatroomRedis/mq"
	"chatroomRedis/routers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	routers.LoadUser(r)				//调用User的路由
	routers.LoadRoom(r)				//调用Room路由
	return r
}

func main(){
	db.InitGorm()				//初始化数据库
	defer db.DB.Close() 		//关闭数据库

	mq.InitRedis()				//初始化redis
	defer mq.Rdb.Close()

	r := SetupRouter()		//设置路由
	r.Run()					//启动服务
}

