package main

import (
	"LRU/db"
	"LRU/router"
	"github.com/gin-gonic/gin"
)

func SetupRouter()*gin.Engine{
	r := gin.New()
	router.LoadBlog(r)		//博客路由
	return r
}

func main(){
	db.InitDB()				//初始化数据库
	defer db.DB.Close()

	r := SetupRouter()
	r.Run(":9090")
}
