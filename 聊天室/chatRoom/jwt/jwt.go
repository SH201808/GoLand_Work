package jwt

import (
	"chatroomRedis/Error"
	"chatroomRedis/db"
	"chatroomRedis/model"
	"chatroomRedis/mq"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type MyClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Minute*10
var MySecret = []byte("chat")

//生成JWT
func GenToken(userId string)string{
	c:=MyClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: 	time.Now().Add(TokenExpireDuration).Unix(),	//token的过期时间
			Issuer:		"zhy",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,c)	// 使用签名方法创建签名对象
	tokenstring,_ := token.SignedString(MySecret)
	err := mq.Rdb.Do("Set",tokenstring,userId).Err()
	if err != nil {
		
	}
	mq.Rdb.Do("expire",tokenstring,c.ExpiresAt-time.Now().Unix())
	return tokenstring					//使用指定的secret签名并获得完整的编码后的字符串token
}

//jJWT认证
func JWTAuthMiddleware()func(c *gin.Context){
	return func(c *gin.Context){
		//客户端携带token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader ==""{
			Err := Error.ErrInfo(0,"请求头中的token为空","登录失败")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		//按空格进行分割判断请求头是否正确
		parts := strings.SplitN(authHeader," ",2)
		if !(len(parts)==2&&parts[0]=="Bearer"){
			Err := Error.ErrInfo(0,"请求头中auth格式有误","登录失败")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		//从redis数据库中判断token是否有效
		userId,err := mq.Rdb.Get(parts[1]).Result()
		if err != nil {
			Err := Error.ErrInfo(0,"无效的token","登录失败")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		var user model.User
		err = db.DB.Where("user_id = ?",userId).First(&user).Error
		if err != nil {
			Err := Error.ErrInfo(0,"用户不存在","用户不存在")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		c.Set("user",user)	//得到userId的值并保存到上下文中
		c.Set("token",parts[1])
		c.Next()
	}
}