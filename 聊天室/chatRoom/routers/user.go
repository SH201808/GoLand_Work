package routers

import (
	"chatroomRedis/Error"
	"chatroomRedis/db"
	"chatroomRedis/jwt"
	"chatroomRedis/model"
	"chatroomRedis/mq"
	"chatroomRedis/verifyCode"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)



func LoadUser(r *gin.Engine){
	r.Use(gin.Logger())
	userGroup := r.Group("/user",Error.Recovery())		//recovery中间件
	{
		userGroup.POST("/register",userRegister)			//用户注册
		userGroup.POST("/login",userLogin)				//用户登录
		userGroup.POST("/forgetPwd",forgetPwd)			//用户忘记密码得到验证码
		userGroup.PUT("/updatePwd",RateLimitMiddleware(3*time.Second,100),updatePwd)		//用户发送验证码修改密码
		userGroup.PUT("/update",jwt.JWTAuthMiddleware(),infoUpdate)		//修改个人信息
		userGroup.DELETE("/exit",jwt.JWTAuthMiddleware(),exitSystem)		//退出系统
	}
}

func RateLimitMiddleware(fillInterval time.Duration,cap int64)func(c *gin.Context){
	bucket := ratelimit.NewBucket(fillInterval,cap)		// 创建指定填充速率和容量大小的令牌桶
	return func(c *gin.Context) {
		var user model.User
		userEmail := c.PostForm("userEmail")
		err := db.DB.Where("user_email = ?",userEmail).First(&user).Error		//查找用户
		errFlag := false
		var Err Error.Err
		if err != nil {
			Err = Error.ErrInfo(0,"数据库中无该邮箱","邮箱错误")
			errFlag = true
		}else {
			flag, _ := mq.Rdb.SIsMember("blackList", user.UserName).Result()	//查看用户是否在黑名单
			if flag == true {
				Err =Error.ErrInfo(0,"用户在黑名单中，无法发送验证码","已在黑名单中，无法发送验证码")
				errFlag = true
			}else {
				count, _ := mq.Rdb.Incr("count" + user.UserName).Result()		//用户发送验证码次数+1
				if count == 1 {
					mq.Rdb.Expire("count"+user.UserName, 1*time.Minute)
				} else if count >= 4 {
					Err = Error.ErrInfo(0,"发送多次拉入黑名单","发送多次拉入黑名单")		//一分钟内超过三次拉入黑名单
					errFlag = true
				}
				if bucket.TakeAvailable(1) < 1 {
					Err = Error.ErrInfo(0,"无令牌","无令牌，请等待")		//令牌不够
					errFlag = true
				}
			}
		}
		if errFlag{
			c.Set("Error",Err)
			c.Abort()
			return
		}else {
			c.Set("user", user)
			c.Next()
		}
	}
}

func userRegister(c *gin.Context){
	user := model.User{}
	c.ShouldBind(&user) //绑定参数
	var users []model.User
	db.DB.Where("user_id = ?",user.UserId).Or("user_email = ?",user.UserEmail).Find(&users)			//查询数据库
	var Err Error.Err
	if len(users)==0 {
		db.DB.Create(&user)		//创建用户
		Err = Error.ErrInfo(0, "", "注册成功")
	}else{
		Err = Error.ErrInfo(0,"用户已存在","用户Id或邮箱已存在，请重新注册")
	}
	c.Set("Error",Err)
	c.Abort()
}

func userLogin(c *gin.Context) {
	user := model.User{}
	c.ShouldBind(&user)
	err := c.ShouldBind(&user) //绑定获取到的参数
	var Err Error.Err
	if err != nil {
		Err = Error.ErrInfo(1,"输入参数无效","输入参数无效,请重新输入")
	}else{
		//判断用户名和密码是否正确
		err = db.DB.Where("user_id = ? AND user_pwd = ?", user.UserId, user.UserPwd).First(&user).Error
		if err != nil {
			Err = Error.ErrInfo(0,"登录失败","id或密码错误")
		}else {
			token := jwt.GenToken(user.UserId)			//得到token
			reply, _ := mq.Rdb.Get(user.UserName).Bytes()		//查看用户加入的房间
			var roomId []string
			json.Unmarshal(reply, &roomId)
			var rooms []model.Room
			for _, v := range roomId {
				room := model.Room{}
				db.DB.Where("room_id = ?", v).First(&room)
				rooms = append(rooms, room)
			}
			tokenMap := make(map[string]interface{})
			userNotice := make(map[string]interface{})
			tokenMap["token"] = token
			userNotice["msg"] = "登录成功 hello " + user.UserName
			userNotice["rooms"] = rooms
			Err = Error.ErrInfo(0, tokenMap, userNotice)
		}
	}
	c.Set("Error",Err)
	c.Abort()
}

func infoUpdate(c *gin.Context){
	i,_ := c.Get("user")
	user := i.(model.User)
	newInfo := model.User{}
	c.ShouldBind(&newInfo)		//绑定新的信息
	var isEmailbind model.User
	err := db.DB.Where("user_Email = ?",newInfo.UserEmail).First(&isEmailbind).Error
	if err == nil {
		Err := Error.ErrInfo(0,"","邮箱已被绑定")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	db.DB.Model(&user).Updates(&newInfo)	//修改信息
	Err := Error.ErrInfo(0,"","修改成功")
	c.Set("Error",Err)
	c.Abort()
	return
}

func exitSystem(c *gin.Context){
	token,_ := c.Get("token")
	err := mq.Rdb.Del(token.(string)).Err()		//删除数据库中的token
	if err != nil {
		Err := Error.ErrInfo(1,"删除token错误","系统错误")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	iUser,_ := c.Get("user")
	user := iUser.(model.User)
	i,ok:=model.ClientRooms.LoadAndDelete(user.UserName)		//得到用户所有的websocket连接
	var Err Error.Err
	if !ok{
		Err = Error.ErrInfo(1,"删除websocket连接错误","系统错误")
	}else {
		Connes := i.([]*model.WebsocketConn)
		for _, v := range Connes {
			v.Conn.Close()		//关闭每个房间的连接
			mq.Publish("system", v.RoomName, "exit"+user.UserName)	//关闭连接所对应的协程
		}
		Err = Error.ErrInfo(0,"","退出系统成功")
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

func forgetPwd(c *gin.Context){
	var user model.User
	c.ShouldBind(&user)			//绑定参数
	err := db.DB.Where("user_email = ?",user.UserEmail).First(&user).Error
	var Err Error.Err
	if err != nil {
		Err = Error.ErrInfo(0,"数据库中无该邮箱","邮箱错误")
	}else {
		VerifyCode := verifyCode.GenVerifyCode(user.UserId)		//获得验证码
		userNotice := make(map[string]interface{})
		userNotice["验证码"] = VerifyCode
		userNotice["notice"] = "十分钟内有效"
		Err = Error.ErrInfo(0,"",userNotice)
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

func updatePwd(c *gin.Context){
	i,_ := c.Get("user")
	user := i.(model.User)
	VerifyCode := c.PostForm("验证码")			//获取用户发送的验证码
	flag := verifyCode.CodeAuth(VerifyCode,user.UserId)		//判断验证码是否有效
	var Err Error.Err
	if flag{
		newPwd := c.PostForm("newPwd")			//获得新的密码
		db.DB.Model(&user).Update("user_pwd",newPwd)		//修改密码
		Err = Error.ErrInfo(0,"","修改成功")
	}else{
		Err = Error.ErrInfo(0,"","验证码错误")
	}
	c.Set("Error",Err)
	c.Abort()
	return
}