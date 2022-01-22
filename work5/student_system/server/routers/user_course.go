package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"student_system/server/Struct"
	"student_system/server/db"
	"student_system/server/jwt"
)

func chooseCourse(c *gin.Context) {
	var user Struct.User
	user.UserId = c.MustGet("userId").(string)		//用户登录后得到用户ID
	db.DB.First(&user)									//查找用户信息
	i ,_:= c.Get("course")							//得到课程信息

	course := i.(Struct.Course)

	record := Struct.UserCourse{
		UserID: user.UserId,
		CourseID:  course.CourseId,
	}
	tCoursePerson := course.CoursePersonSum + 1
	tUserCredit := user.UserCredit + course.CourseCredit
	err := db.DB.Where("user_id = ? AND course_id = ?", record.UserID, record.CourseID).Find(&record).Error
	//判断记录是否存在且判断课程人数和用户学分是否满
	if err != nil && tCoursePerson <= 2 && tUserCredit <= 10 {
		db.DB.Create(&record)						//创建记录
		course.CoursePersonSum = tCoursePerson
		db.DB.Save(&course)							//更新课程人数

		db.DB.Model(&user).Update("userCredit",tUserCredit)	//更新用户学分
		c.JSON(http.StatusOK, gin.H{
			"msg": "选课成功",
		})
	} else if err == nil {				//课程已选
		c.JSON(http.StatusOK, gin.H{
			"msg": "已选课成功，请勿重复添加",
		})
		return
	} else if tCoursePerson > 2 {		//课程人数已满
		c.JSON(http.StatusOK, gin.H{
			"msg": "该课人数已满，请选择其他课",
		})
		return
	} else if tUserCredit > 10 {		//学分已满
		c.JSON(http.StatusOK, gin.H{
			"msg": "学分已满",
		})
		return
	}
}

func deleteUserCourse(c *gin.Context) {
	var user Struct.User
	user.UserId = c.MustGet("userId").(string)		//得到用户登录后用户的ID
	db.DB.First(&user)									//得到用户信息
	i, _ := c.Get("course")						//得到课程信息
	course := i.(Struct.Course)
	var record Struct.UserCourse
	//判断是否有该记录
	err := db.DB.Where("user_id = ? AND course_id = ?", user.UserId, course.CourseId).Find(&record).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "记录不存在",
		})
		return
	}
	db.DB.Delete(&record)		//删除记录
	course.CoursePersonSum--	//课程人数减1
	db.DB.Save(&course)			//更新课程数据
	user.UserCredit -= course.CourseCredit			//用户学分减少
	db.DB.Save(&user)			//更新用户信息
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func  RetrieveUserCourse(c *gin.Context)  {
	var userCourse Struct.UserCourse
	err := c.ShouldBind(&userCourse)		//得到查找的参数
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"参数有误",
		})
		return
	}
	//查找记录
	err = db.DB.Where("user_id= ? And course_id = ?",userCourse.UserID,userCourse.CourseID).First(&userCourse).Error
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"记录不存在",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"选课记录":userCourse,		//输出记录
		"msg":"查询完成",
	})
}

func LoadUserCourse(e *gin.Engine){
	userCourse := e.Group("/user_course")
	{
		userCourse.GET("/retrieve", RetrieveUserCourse) 										//user_courses表的检索
		userCourse.DELETE("/delete",jwt.JWTAuthMiddleware(),FindCourse(),deleteUserCourse)	//user_courses表的删除
		userCourse.POST("/create",jwt.JWTAuthMiddleware(),FindCourse(),chooseCourse)			//user_courses表的增添
	}
}