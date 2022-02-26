package db

import (
	"chatroomRedis/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
	err error
)

func InitGorm(){
	//连接数据库
	DB,err = gorm.Open("mysql","root:1358@(127.0.0.1:3306)/chatroom?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	sqlDB := DB.DB()
	sqlDB.SetMaxIdleConns(50)		//设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)		//设置打开数据库连接的最大数量
	//自动迁移
	DB.AutoMigrate(&model.User{},&model.Room{})

}