package Error

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Err struct {
	Code int
	Msg interface{}
	Notice interface{}
}

func ErrInfo(code int,msg,notice interface{})Err{
	err := Err{
		code,
		msg,
		notice,
	}
	return err
}

func Recovery()gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err, _ := c.Get("Error")			//得到返回的错误
		c.JSON(http.StatusOK, gin.H{
			"err": err.(Err),
		})
		c.Abort()
		return
	}
}