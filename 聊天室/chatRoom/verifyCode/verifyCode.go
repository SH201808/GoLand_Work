package verifyCode

import (
	"chatroomRedis/mq"
	"fmt"
	"math/rand"
	"time"
)

//生成验证码
func GenVerifyCode(userId string)string{
	//生成四位验证码
	codeString := fmt.Sprintf("%04v",rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000))
	mq.Rdb.Set("code"+userId,codeString,10*time.Minute)		//将验证码放入redis中，设置十分钟内有效
	return codeString
}

//判断验证码是否有效
func CodeAuth(verifyCode string,userId string)bool{
	InDbCode,err := mq.Rdb.Get("code"+userId).Result()		//从redis中获取验证码
	//判断验证码是否有效
	if err != nil ||InDbCode!=verifyCode{
		return false
	}
	return true
}
