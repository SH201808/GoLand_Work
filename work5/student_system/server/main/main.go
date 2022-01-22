package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"student_system/server/db"
	"student_system/server/routers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	routers.LoadUser(r)				//调用User的路由
	routers.LoadCourse(r)			//调用Course的路由
	routers.LoadUserCourse(r)		//调用User_Course的路由

	return r
}

func main(){
	db.InitDB()				//初始化数据库
	defer db.DB.Close()		//关闭数据库

	r := SetupRouter()		//设置路由
	r.Run()					//启动服务
}


