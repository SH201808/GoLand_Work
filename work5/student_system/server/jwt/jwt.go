package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour*2
var MySecret = []byte("Student_System")

//生成JWT
func GenToken(userId string)(string,error){
	c:=MyClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: 	time.Now().Add(TokenExpireDuration).Unix(),	//token的过期时间
			Issuer:		"zhy",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,c)	// 使用签名方法创建签名对象
	return token.SignedString(MySecret)						//使用指定的secret签名并获得完整的编码后的字符串token
}

//解析JWT
func ParseToken(tokenString string)(*MyClaims,error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	//校验token
	claims, ok := token.Claims.(*MyClaims)
	if ok && token.Valid {
		return claims,nil
	}
	return nil,errors.New("invalid token")
}

//
func JWTAuthMiddleware()func(c *gin.Context){
	return func(c *gin.Context){
		//客户端携带token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader ==""{
			c.JSON(http.StatusOK,gin.H{
				"msg":"请求头中的auth为空",
			})
			c.Abort()
			return
		}
		//按空格进行分割判断请求头是否正确
		parts := strings.SplitN(authHeader," ",2)
		if !(len(parts)==2&&parts[0]=="Bearer"){
			c.JSON(http.StatusOK,gin.H{
				"msg":"请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		mc,err:=ParseToken(parts[1])	//得到tokenstring并解析
		if err != nil {
			c.JSON(http.StatusOK,gin.H{
				"msg":"无效的Token",
			})
			c.Abort()
			return
		}
		c.Set("userId",mc.UserId)	//得到userId的值并保存到上下文中
		c.Next()
	}
}