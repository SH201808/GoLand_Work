package db

import (
	"github.com/jinzhu/gorm"
	"student_system/server/Struct"
)

var (
	DB *gorm.DB
	err error
)

func InitDB(){
	//连接数据库
	DB,err = gorm.Open("mysql","root:1358@(127.0.0.1:3306)/student_system?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	sqlDB := DB.DB()
	sqlDB.SetMaxIdleConns(50)		//设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)		//设置打开数据库连接的最大数量
	//自动迁移并设置外键
	DB.AutoMigrate(&Struct.User{},&Struct.Course{},&Struct.UserCourse{})
	DB.Model(&Struct.UserCourse{}).AddForeignKey("user_id","users(user_id)","CASCADE","CASCADE")
	DB.Model(&Struct.UserCourse{}).AddForeignKey("course_id","courses(course_id)","CASCADE","CASCADE")
}