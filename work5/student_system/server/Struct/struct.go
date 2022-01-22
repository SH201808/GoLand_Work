package Struct

type User struct{
	UserId string `gorm:"PRIMARY_KEY" form:"userId"`
	UserPwd string `form:"userPwd"`
	UserName string	`form:"userName"`
	UserCredit int8
}

type Course struct {
	CourseId string `gorm:"PRIMARY_KEY;AUTO_INCREMENT" form:"courseId"`
	CourseName string	`form:"courseName"`
	CourseCredit int8	`form:"courseCredit"`
	CoursePersonSum int8
}

type UserCourse struct {
	ID int `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserID	string	`form:"userId" binding:"required"`
	CourseID string	`form:"courseId" binding:"required"`
}
