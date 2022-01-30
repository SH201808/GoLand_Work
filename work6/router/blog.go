package router

import (
	"LRU/client"
	"LRU/db"
	"LRU/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadBlog(r *gin.Engine){
	r.Use(recovery())			//设置中间件统一返回错误
	blogGroup := r.Group("/blog")
	{
		blogGroup.POST("/create",createRecord)		//新增博客数据
		blogGroup.GET("/retrieve",retrieveRecord)	//查询博客数据
		blogGroup.PUT("/update",updateRecord)		//更新博客数据
		blogGroup.DELETE("/delete",deleteRecord)		//删除博客数据
	}
}

func recovery()gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err, _ := c.Get("err")			//得到返回的错误
		c.JSON(http.StatusOK, gin.H{
			"err": err.(utils.Err),
		})
		NewErr, ok := c.Get("NewErr")
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"err": NewErr.(utils.Err),
			})
		}
		return
	}
}

func createRecord(c *gin.Context) {
	blog := db.Blog{}
	err := c.ShouldBind(&blog)		//绑定参数
	if err != nil {
		Err := utils.ErrInfo(1,"参数有误","请按规范输入参数")
		Err.Error()
		c.Set("err",Err)
		c.Abort()
		return
	}
	err = db.DB.First(&blog).Error	//在数据库中判断记录是否存在
	if err != nil {
		db.DB.Create(&blog)
	}
	client := client.NewClient("127.0.0.1:8080")	//开启一个客户端
	err = client.Set(blog.Id, blog.Title)	//更新缓存
	var Err utils.Err
	if err == nil {
		Err = utils.ErrInfo(0,"","创建成功")
	}else {
		Err = utils.ErrInfo(1, "发送或接受数据有误", "请重新发送数据")
	}
	Err.Error()
	c.Set("err",Err)
	c.Abort()
}

func retrieveRecord(c *gin.Context) {
	blog := db.Blog{}
	blog.Id = c.Query("Id")		//得到参数
	client := client.NewClient("127.0.0.1:8080")		//开启一个新的客户端
	v, err := client.Get(blog.Id)				//在缓存中查询数据
	if err != nil {
		Err := utils.ErrInfo(1, "查询缓存发生错误", "请重新发送数据")
		Err.Error()
		c.Set("err",Err)
		c.Abort()
		return
	}
	var Err utils.Err
	if v == "LRU中无记录" {
		err = db.DB.First(&blog).Error		//在数据库中查找数据
		if err != nil {
			Err = utils.ErrInfo(0,"","数据库中无该记录")
		}else{
			Err = utils.ErrInfo(0,"",blog.Title)
		}
		err = client.Set(blog.Id, blog.Title)	//更新缓存
		if err != nil {
			NewErr := utils.ErrInfo(1,"更新缓存发生错误","请重新发送")
			c.Set("NewErr",NewErr)
			c.Abort()
		}
	} else {
		Err = utils.ErrInfo(0,"",v)
	}
	c.Set("err",Err)
	c.Abort()
}

func updateRecord(c *gin.Context){
	blog := db.Blog{}
	blog.Id = c.Query("Id")	//绑定参数
	err := db.DB.First(&blog).Error	//数据库中查询记录
	if err != nil {
		Err := utils.ErrInfo(0,"","数据库中无记录无法修改")
		c.Set("err",Err)
		c.Abort()
		return
	}
	tempBlog := db.Blog{}
	tempBlog.Title = c.PostForm("title")
	db.DB.Model(&blog).Updates(tempBlog)	//修改数据库中的数据
	client := client.NewClient("127.0.0.1:8080")
	err = client.Delete(blog.Id)			//删除缓存
	var Err utils.Err
	if err != nil {
		Err = utils.ErrInfo(1,"删除缓存发送错误","请重新更新")
	}else{
		Err = utils.ErrInfo(0,"","修改成功")
	}
	c.Set("err",Err)
	c.Abort()
}

func deleteRecord(c *gin.Context){
	blog := db.Blog{}
	blog.Id = c.Query("Id")
	err := db.DB.First(&blog).Error		//数据库中查询数据
	if err != nil {
		Err := utils.ErrInfo(0,"","数据库中无记录无法删除")
		c.Set("err",Err)
		c.Abort()
		return
	}
	db.DB.Delete(&blog)		//删除数据
	client := client.NewClient("127.0.0.1:8080")
	err = client.Delete(blog.Id)	//删除缓存
	var Err utils.Err
	if err != nil {
		Err = utils.ErrInfo(1,"删除缓存发送错误","请重新更新")
	}else{
		Err = utils.ErrInfo(0,"","删除成功")
	}
	c.Set("err",Err)
	c.Abort()
}