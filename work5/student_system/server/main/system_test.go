package main

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"student_system/server/db"
	"testing"
)

func Test_userLogin(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"登录成功", "userId=1024&userPwd=1358", "登录成功 hello test"},
		{"登录失败", "userId=255&userPwd=1358", "登录失败"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/user/login", strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_userRegister(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"注册成功", "userId=1024&userPwd=1358&userName=test", "注册成功"},
		{"用户已存在", "userId=1024&userPwd=1358&userName=test", "用户已存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/user/register", strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_userDelete(t *testing.T) {
	tests:= []struct{
		name string
		params string
		except string
	}{
		{"删除用户","userId=1024","删除成功"},
		{"用户不存在","userId=1024","用户不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/user/delete?" + tt.params
			req := httptest.NewRequest("DELETE", url, strings.NewReader(tt.params))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_userRetrieve(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"查询用户信息", "userId=1024", "查询完成"},
		{"用户不存在", "userId=1023", "用户不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/user/retrieve?"+tt.params
			req := httptest.NewRequest("GET", url, strings.NewReader(tt.params))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_userUpdate(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"修改用户信息", "userId=1024", "修改成功"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/user/update?"+tt.params
			p := "userPwd=135"
			req := httptest.NewRequest("PUT", url, strings.NewReader(p))
			req.Header.Set("Content-Type","application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_courseCreate(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"创建课程", "courseId=1024&courseName=test&courseCredit=1", "创建课程成功"},
		{"课程已存在", "courseId=1024&courseName=test&courseCredit=1", "课程已存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/course/create", strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_courseRetrieve(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"查询课程信息", "courseId=1024", "查询完成"},
		{"课程不存在", "courseId=1023", "课程不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/course/retrieve?"+tt.params
			req := httptest.NewRequest("GET", url, strings.NewReader(tt.params))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_courseUpdate(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"修改课程信息", "courseId=1024", "修改成功"},
		{"课程不存在", "courseId=1023", "课程不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/course/update?"+tt.params
			p := "courseName=tes"
			req := httptest.NewRequest("PUT", url, strings.NewReader(p))
			req.Header.Set("Content-Type","application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_courseDelete(t *testing.T) {
	tests:= []struct{
		name string
		params string
		except string
	}{
		{"删除课程","courseId=1024","删除成功"},
		{"课程不存在","courseId=1024","课程不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/course/delete?" + tt.params
			req := httptest.NewRequest("DELETE", url, strings.NewReader(tt.params))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_userChooseCourse(t *testing.T) {
	tests:= []struct{
		name string
		params string
		except string
	}{
		{"选择课程","courseId=1","选课成功"},
		{"重复选择","courseId=1","已选课成功，请勿重复添加"},
		{"学分已满","courseId=2","学分已满"},
		{"课程人数已满","courseId=1024","该课人数已满，请选择其他课"},
		{"课程不存在","courseId=255","课程不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	reqLogin := httptest.NewRequest("POST","/user/login",strings.NewReader("userId=1024&userPwd=1358"))
	reqLogin.Header.Set("Content-Type","application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w,reqLogin)
	var token map[string]string
	err := json.Unmarshal([]byte(w.Body.String()),&token)
	assert2.Nil(t, err)

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url :="/user_course/create?"+tt.params
			req := httptest.NewRequest("POST", url, strings.NewReader(tt.params))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Authorization","Bearer "+token["token"])
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err = json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_deleteUserCourse(t *testing.T) {
	tests:= []struct{
		name string
		params string
		except string
	}{
		{"删除选课记录","courseId=1","删除成功"},
		{"重复删除","courseId=1","记录不存在"},
		{"课程不存在","courseId=3","课程不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	reqLogin := httptest.NewRequest("POST","/user/login",strings.NewReader("userId=1024&userPwd=1358"))
	reqLogin.Header.Set("Content-Type","application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w,reqLogin)
	var token map[string]string
	err := json.Unmarshal([]byte(w.Body.String()),&token)
	assert2.Nil(t, err)

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/user_course/delete?" + tt.params
			req := httptest.NewRequest("DELETE", url, strings.NewReader(tt.params))
			req.Header.Set("Authorization","Bearer "+token["token"])
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err = json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}

func Test_userCourseRetrieve(t *testing.T) {
	tests := []struct{
		name string
		params string
		except string
	}{
		{"查询选课记录", "userId=1024&courseId=1", "查询完成"},
		{"选课记录不存在", "userId=1023&courseId=1", "记录不存在"},
	}
	r := SetupRouter()
	db.InitDB()
	defer db.DB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/user_course/retrieve?"+tt.params
			req := httptest.NewRequest("GET", url, strings.NewReader(tt.params))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert2.Nil(t, err)
			assert.Equal(t, tt.except, resp["msg"])
		})
	}
}
