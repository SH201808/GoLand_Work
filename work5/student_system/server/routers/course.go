package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"student_system/server/Struct"
	"student_system/server/db"
)

func FindCourse()gin.HandlerFunc {
	return func(c *gin.Context) {
		course := Struct.Course{
			CourseId: c.Query("courseId"),			//获得参数
			CourseName: c.Query("courseName"),
		}
		err := db.DB.Where(course).First(&course).Error	//查询是否有该课
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "课程不存在",
			})
			c.Abort()
			return
		}
		c.Set("course", course)					//将查到的数据进行set，实现跨中间件取值
		c.Next()
	}
}

func CreateCourse(c *gin.Context) {
	var course Struct.Course
	err := c.ShouldBind(&course)		//得到参数
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"参数有误",
		})
		return
	}
	err = db.DB.First(&course).Error	//查询是否有该课
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "课程已存在",
		})
		return
	}
	course.CoursePersonSum = 0
	db.DB.Create(&course)			//创建课程
	c.JSON(http.StatusOK,gin.H{
		"msg":"创建课程成功",
	})
}

func LoadCourse(e *gin.Engine){
	CourseGroup := e.Group("/course")
	{

		//courses表的创建
		CourseGroup.POST("/create", CreateCourse)

		//courses表的检索
		CourseGroup.GET("/retrieve", FindCourse(), func(c *gin.Context) {
			i,_ := c.Get("course")				//得到查询的数据
			course := i.(Struct.Course)
			c.JSON(http.StatusOK, gin.H{
				"课程": course,						//输出数据
				"msg":"查询完成",
			})
		})

		//courses表的更新
		CourseGroup.PUT("/update",FindCourse(),func(c *gin.Context) {
			var course,tCourse Struct.Course
			course.CourseId = c.Query("courseId")	//得到需要更新的课程ID
			err := c.ShouldBind(&tCourse)				//绑定更改的值
			if err != nil {
				c.JSON(http.StatusOK,gin.H{
					"msg":"参数有误",
				})
				return
			}
			db.DB.Model(&course).Updates(&tCourse)		//更新数据
			c.JSON(http.StatusOK,gin.H{
				"msg":"修改成功",
			})
		})

		//courses表的删除
		CourseGroup.DELETE("/delete", FindCourse(), func(c *gin.Context) {
			i,_ :=c.Get("course")				//得到检索后的数据
			course := i.(Struct.Course)
			db.DB.Delete(&course)					//删除数据
			c.JSON(http.StatusOK,gin.H{
				"msg":"删除成功",
			})
		})

	}
}