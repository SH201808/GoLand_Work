package main


//未完成
import (
	"chatroomRedis/db"
	"chatroomRedis/mq"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_roomCreate (t *testing.T){
	tests := []struct{
		name string
		params string
		except string
	}{
		{"创建成功", "roomName=go&roomId=1&roomCap=3&roomAccess=1", "修改成功"},
	}
	r := SetupRouter()
	db.InitGorm()
	defer db.DB.Close()
	mq.InitRedis()
	defer mq.Rdb.Close()

	for _,tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			token := requestLogin(r,"userId=1024&userPwd=135")

			url := "/room/create?"+tt.params
			req := httptest.NewRequest("GET", url, strings.NewReader(tt.params))
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