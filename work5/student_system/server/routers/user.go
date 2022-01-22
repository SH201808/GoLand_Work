package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"student_system/server/Struct"
	"student_system/server/db"
	"student_system/server/jwt"
)



func FindUser()gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []Struct.User
		user := Struct.User{}
		user.UserId = c.Query("userId")		//获得参数
		user.UserName = c.Query("userName")
		db.DB.Where(user).Find(&users)				//找到符合参数的user信息
		if len(users) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "用户不存在",
			})
			c.Abort()
			return
		} else if len(users) < 2 {
			c.Set("user", users[0])		//只有一个user符合，可以进行删除或者修改
		}
		c.Set("users", users)
		c.Next()
	}
}

func UserRegister(c *gin.Context){
	var user Struct.User
	err := c.ShouldBind(&user)			//得到参数
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"参数有误",
		})
		return
	}
	err = db.DB.First(&user).Error		//寻找是否有符合参数的user
	if err == nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"用户已存在",
		})
		return
	}
	user.UserCredit =0
	db.DB.Create(&user)					//创建user
	c.JSON(http.StatusOK,gin.H{
		"msg":"注册成功",
	})
}

func UserLogin(c *gin.Context) {
	var user Struct.User
	err := c.ShouldBind(&user)		//绑定获取到的参数
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "输入参数无效",
		})
		return
	}
	//判断用户名和密码是否正确
	err = db.DB.Where("user_id = ? AND user_pwd = ?", user.UserId, user.UserPwd).First(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录失败",
		})
		return
	}
	tokenString, _ := jwt.GenToken(user.UserId)		//得到用户的token
	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功 hello " + user.UserName,
		"token": tokenString,	//返回token
	})
}

func LoadUser(e *gin.Engine){
	userGroup := e.Group("/user")
	{
		//user的登录
		userGroup.POST("/login", UserLogin)

		//检索user的信息
		userGroup.GET("/retrieve",FindUser(), func(c *gin.Context) {
			i, _ := c.Get("users")			//从上下文中得到users切片
			users := i.([]Struct.User)
			c.JSON(http.StatusOK, gin.H{
				"用户信息":users,
				"msg":"查询完成",
			})
		})

		//users表的增添
		userGroup.POST("/register",UserRegister)

		//users表的删除
		userGroup.DELETE("/delete", FindUser(), func(c *gin.Context) {
			i, _ := c.Get("user")	//得到检索后符合条件的user
			user := i.(Struct.User)
			db.DB.Delete(&user)			//删除user
			c.JSON(http.StatusOK, gin.H{
				"msg": "删除成功",
			})
		})

		//users表的更新
		userGroup.PUT("/update",FindUser(), func(c *gin.Context) {
			var user,tUser Struct.User
			i,_ := c.Get("user")		//得到符合条件的user
			user = i.(Struct.User)
			err := c.ShouldBind(&tUser)		//得到修改的参数
			if err != nil {
				c.JSON(http.StatusOK,gin.H{
					"msg":"参数有误",
				})
				return
			}
			db.DB.Model(&user).Updates(&tUser)	//根据新参数对user信息进行修改
			c.JSON(http.StatusOK,gin.H{
				"msg":"修改成功",
			})
		})
	}
}