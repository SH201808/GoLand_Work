package main

import (
	"chatroomRedis/db"
	"chatroomRedis/mq"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_userRegister(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"注册成功", "userId=1024&userPwd=1358&userName=test", "注册成功"},
		{"用户已存在", "userId=1024&userPwd=1358&userName=test", "用户已存在，请重新注册"},
	}
	r := SetupRouter()			//设置引擎
	db.InitGorm()				//初始化数据库
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//发起请求
			req := httptest.NewRequest("POST", "/user/register", strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)		//接受响应
			assert2.Nil(t, err)
			noticeMap := resp["err"]
			var notice string
			for k,v := range noticeMap.(map[string]interface{}){
				if k == "Notice"{
					notice = v.(string)
					break
				}
			}
			assert.Equal(t,tt.except,notice)
		})
	}
}

func Test_userLogin(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"登录成功", "userId=1024&userPwd=1358", "登录成功 hello test"},
		{"登录失败", "userId=255&userPwd=1358", "id或密码错误"},
	}
	r := SetupRouter()		//设置引擎
	db.InitGorm()			//初始化mysql
	defer db.DB.Close()
	mq.InitRedis()			//初始化redis
	defer mq.Rdb.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//发起请求
			req := httptest.NewRequest("POST", "/user/login", strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)		//接受响应
			assert2.Nil(t, err)
			noticeMap := resp["err"]
			var notice string
			for k1,v1 := range noticeMap.(map[string]interface{}){
				if k1 == "Notice" {
					if _,ok := v1.(string);ok{
						notice = v1.(string)
						break
					}else {
						for k2, v2 := range v1.(map[string]interface{}) {
							if k2 == "msg" {
								notice = v2.(string)
								break
							}
						}
					}
				}
			}
			assert.Equal(t,tt.except,notice)
		})
	}
}


func Test_forgetPwd_updatePwd(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"发送验证码成功", "userEmail=1@qq.com", "十分钟内有效"},
		{"发送验证码失败", "userEmail=qq.com", "邮箱错误"},
	}
	r := SetupRouter()
	db.InitGorm()
	defer db.DB.Close()
	mq.InitRedis()
	defer mq.Rdb.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//发起请求
			req := httptest.NewRequest("POST", "/user/forgetPwd", strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)		//接受响应
			assert2.Nil(t, err)
			noticeMap := resp["err"]
			var notice string
			for k1,v1 := range noticeMap.(map[string]interface{}){
				if k1 == "Notice" {
					if _,ok := v1.(string);ok{
						notice = v1.(string)
						break
					}else {
						for k2, v2 := range v1.(map[string]interface{}) {
							if k2 == "notice" {
								notice = v2.(string)
								break
							}
						}
					}
				}
			}
			assert.Equal(t,tt.except,notice)
		})
	}
}

func Test_infoUpdate(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"修改成功", "userId=1024&userPwd=1358", "修改成功"},
	}
	r := SetupRouter()
	db.InitGorm()
	defer db.DB.Close()
	mq.InitRedis()
	defer mq.Rdb.Close()


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//发起登录请求
			token := requestLogin(r,tt.params)

			//发起修改信息请求
			req := httptest.NewRequest("PUT", "/user/update", strings.NewReader("userId=1024&userPwd=135"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization","Bearer "+token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			noticeMap := resp["err"]
			var notice string
			for k1,v1 := range noticeMap.(map[string]interface{}){
				if k1 == "Notice" {
					if _,ok := v1.(string);ok{
						notice = v1.(string)
						break
					}else {
						for k2, v2 := range v1.(map[string]interface{}) {
							if k2 == "msg" {
								notice = v2.(string)
								break
							}
						}
					}
				}
			}
			assert.Equal(t,tt.except,notice)
		})
	}
}


func requestLogin(r *gin.Engine,params string)(token string){
	reqLogin := httptest.NewRequest("POST","/user/login",strings.NewReader(params))
	reqLogin.Header.Set("Content-Type","application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w,reqLogin)
	var resp map[string]interface{}
	json.Unmarshal([]byte(w.Body.String()),&resp)		//获取token
	tokenMap := resp["err"]
	for k,v := range tokenMap.(map[string]interface{}){
		if k == "Msg"{
			for k2,v2 := range v.(map[string]interface{}){
				if k2 == "token"{
					token = v2.(string)
				}
			}
		}
	}
	return token
}