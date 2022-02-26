package routers

import (
	"bytes"
	"chatroomRedis/Error"
	"chatroomRedis/db"
	"chatroomRedis/jwt"
	"chatroomRedis/model"
	"chatroomRedis/mq"
	"chatroomRedis/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)



func LoadRoom(r *gin.Engine){
	roomGroup := r.Group("/room")
	{
		roomGroup.GET("/create",jwt.JWTAuthMiddleware(),roomCreate)
		roomGroup.GET("/enter",jwt.JWTAuthMiddleware(),enterRoom)
		roomGroup.GET("/reEnter",jwt.JWTAuthMiddleware(),reEnterRoom)
		roomGroup.GET("/retrieveRoom", Error.Recovery(),jwt.JWTAuthMiddleware(),roomRetrieve)
		roomGroup.POST("/sendImg", Error.Recovery(),jwt.JWTAuthMiddleware(),sendImg)
		roomGroup.DELETE("/exit", Error.Recovery(),jwt.JWTAuthMiddleware(),exitRoom)
		roomGroup.GET("/retrieveUser", Error.Recovery(),jwt.JWTAuthMiddleware(),userRetrieve)
		roomGroup.DELETE("/kickOutUser", Error.Recovery(),jwt.JWTAuthMiddleware(),OwnerAuth(),kickOutUser)
		roomGroup.PUT("/transferOwner", Error.Recovery(),jwt.JWTAuthMiddleware(),OwnerAuth(),transferOwner)
		roomGroup.PUT("/update", Error.Recovery(),jwt.JWTAuthMiddleware(),OwnerAuth(),updateRoom)
		roomGroup.DELETE("/dissolve", Error.Recovery(),jwt.JWTAuthMiddleware(),OwnerAuth(),dissolveRoom)
	}
}

func OwnerAuth()gin.HandlerFunc{
	return func(c *gin.Context) {
		i, _ := c.Get("user")		//得到认证通过的用户
		user := i.(model.User)
		room := model.Room{}
		room.RoomId = c.Query("roomId")		//得到roomId
		err := db.DB.Where("room_id = ?",room.RoomId).First(&room).Error		//查询房间
		if err != nil {
			Err := Error.ErrInfo(0,"房间不存在","房间不存在")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		//判断是否有权限执行踢出用户等操作
		if room.RoomOwner != user.UserId{
			Err := Error.ErrInfo(0,"不是房主无权限","你不是房主无权限")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		c.Set("room",room)
	}
}

func roomCreate(c *gin.Context) {
	room := model.Room{}
	err := c.ShouldBind(&room)		//得到房间参数
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.DB.Where("room_id = ?",room.RoomId).First(&room).Error		//查询是否有房间
	if err == nil {
		return
	}
	room.CreatedAt = time.Now()			//创建时间
	i, _ := c.Get("user")
	user := i.(model.User)
	room.RoomOwner = user.UserId		//房主
	db.DB.Create(&room)					//创建房间
	//升级ws协议
	wsconn, err := (&websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	client := model.Client{
		UserId:   user.UserId,
		UserName: user.UserName,
	}
	clients := make([]model.Client,0)
	clients = append(clients,client)
	data,err := json.Marshal(clients)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = mq.Rdb.Do("Set",room.RoomId,data).Err()		//储存房间内所有的用户信息
	if err != nil {
		fmt.Println(err)
		return
	}
	//将新房间添加到用户加入的所有房间中
	reply,err := mq.Rdb.Get(user.UserName).Bytes()
	if err != nil &&err.Error()=="redis nil" {
		fmt.Println(err)
		return
	}
	var rooms []string
	json.Unmarshal(reply,&rooms)
	rooms = append(rooms,room.RoomId)
	data,err = json.Marshal(rooms)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = mq.Rdb.Do("Set",user.UserName,data).Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.EnterRoom(user.UserName,room.RoomName,wsconn)			//进入房间
}

func roomRetrieve(c *gin.Context){
	room := model.Room{}
	err := c.ShouldBind(&room)			//绑定房间信息
	if err != nil {
		Err := Error.ErrInfo(0,"参数无效","参数无效")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	var rooms []model.Room
	//按房间名查找并按创建时间排序
	err = db.DB.Where("room_name LIKE ?","%"+room.RoomName+"%").Order("created_at desc").Find(&rooms).Error
	var Err Error.Err
	if err != nil {
		Err = Error.ErrInfo(1,"","未找到房间")
	}else {
		userNotice := make(map[string]interface{})
		userNotice["rooms"] = rooms		//搜索到的房间列表
		Err = Error.ErrInfo(1,"",userNotice)
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

func enterRoom(c *gin.Context){
	room := model.Room{}
	c.ShouldBind(&room)				//绑定参数
	err := db.DB.Where("room_id = ?",room.RoomId).First(&room).Error		//查找房间
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"房间未找到",
		})
		return
	}
	//判断房间权限
	if !room.RoomAccess{
		c.JSON(http.StatusOK,gin.H{
			"msg":"房间不允许其他人进入",
		})
		return
	}
	i,_ := c.Get("user")
	user := i.(model.User)
	flag,clients,_,userNumber,err := utils.UserInRoom(room.RoomId,user.UserId)			//判断用户是否在房间内
	if err != nil {
		fmt.Println(err)
		return
	}
	if flag{
		c.JSON(http.StatusOK,gin.H{
			"msg":"已在房间内",
		})
		return
	}
	//判断房间人数是否已满
	if userNumber == room.RoomCap{
		c.JSON(http.StatusOK,gin.H{
			"msg":"房间人数已满",
		})
		return
	}
	client := model.Client{
		UserId:   user.UserId,
		UserName: user.UserName,
	}
	clients = append(clients,client)
	data,err := json.Marshal(clients)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = mq.Rdb.Do("Set",room.RoomId,data).Err()		//添加房间人数
	if err != nil {
		fmt.Println(err)
		return
	}
	//将新加入的房间添加到用户下
	reply,err := mq.Rdb.Get(user.UserName).Bytes()
	if err != nil && err.Error()!="redis: nil" {
		fmt.Println(err)
		return
	}
	var rooms []string
	json.Unmarshal(reply,&rooms)
	rooms = append(rooms,room.RoomId)
	data,err = json.Marshal(rooms)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = mq.Rdb.Do("Set",user.UserName,data).Err()
	if err != nil {
		fmt.Println(err)
		return
	}
	//升级ws协议
	wsconn,err := (&websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer,c.Request,nil)
	utils.EnterRoom(user.UserName,room.RoomName,wsconn)			//进入房间
}

//重新登陆后进入已经进入的房间
func reEnterRoom(c *gin.Context){
	room := model.Room{}
	c.ShouldBind(&room)			//绑定参数
	err := db.DB.Where("room_id = ?",room.RoomId).First(&room).Error		//查询房间
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"房间未找到",
		})
		return
	}
	i,_ := c.Get("user")
	user := i.(model.User)
	//升级ws协议
	wsconn,err := (&websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer,c.Request,nil)
	utils.EnterRoom(user.UserName,room.RoomName,wsconn)		//进入房间
}

func exitRoom(c *gin.Context) {
	room := model.Room{}
	c.ShouldBind(&room)				//绑定参数
	err := db.DB.Where("room_id = ?", room.RoomId).First(&room).Error		//查找房间
	if err != nil {
		Err := Error.ErrInfo(0,"","房间不存在")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	i, _ := c.Get("user")
	user := i.(model.User)
	flag,clients,index,_,err := utils.UserInRoom(room.RoomId,user.UserId)		//判断用户是否在房间内
	var Err Error.Err
	if flag {
		if err != nil {
			Err = Error.ErrInfo(1,err,"退出房间错误")
			c.Set("Error",Err)
			c.Abort()
			return
		}
		err = utils.ExitRoom(clients,index,room.RoomId,room.RoomName,user.UserId,user.UserName)		//退出房间
		if err != nil {
			Err = Error.ErrInfo(1,err,"退出房间错误")
		}else {
			Err = Error.ErrInfo(0,"","退出房间成功")
		}
		c.Set("Error",Err)
		c.Abort()
		return
	}else{
		Err = Error.ErrInfo(0,"","用户不在房间内")
		c.Set("Error",Err)
		c.Abort()
		return
	}
}

func userRetrieve(c *gin.Context){
	room := model.Room{}
	c.ShouldBind(&room)			//绑定参数
	err := db.DB.Where("room_id = ?",room.RoomId).First(&room).Error		//查找房间
	if err != nil {
		Err := Error.ErrInfo(0,"","房间不存在")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	i,_ := c.Get("user")
	user := i.(model.User)
	flag,clients,_,_,err:=utils.UserInRoom(room.RoomId,user.UserId)			//判断用户是否在房间内
	if err != nil {
		Err := Error.ErrInfo(1,err,"查找错误")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	var Err Error.Err
	if !flag {
		Err = Error.ErrInfo(0,"","用户不在房间内，无权限查找")
	}else{
		userNotice := make(map[string]interface{})
		userNotice["clients"] = clients			//用户列表
		Err = Error.ErrInfo(0,"",userNotice)
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

//踢出用户
func kickOutUser(c *gin.Context) {
	userId := c.Query("userId")		//用户参数
	r, _ := c.Get("room")				//房间参数
	room := r.(model.Room)
	flag, clients, index,_, err := utils.UserInRoom(room.RoomId, userId)		//判断用户是否在房间内
	if err != nil {
		Err := Error.ErrInfo(1,err,"踢出错误")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	var Err Error.Err
	if !flag {
		Err = Error.ErrInfo(0,"","用户不在房间内")
	}else{
		var user model.User
		db.DB.Where("user_id = ?",userId).First(&user)
		err = utils.ExitRoom(clients,index,room.RoomId,room.RoomName,user.UserId,user.UserName)		//踢出房间
		if err != nil {
			Err = Error.ErrInfo(1,err,"踢出房间错误")
		}else{
			Err = Error.ErrInfo(0,"","踢出成功")
		}
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

func transferOwner(c *gin.Context) {
	newOwener := c.Query("userId")		//新房主ID
	r,_ := c.Get("room")				//房间参数
	room := r.(model.Room)
	flag,_,_,_,_:=utils.UserInRoom(room.RoomId,newOwener)
	if !flag{
		Err := Error.ErrInfo(0,"","用户不在房间内")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	room.RoomOwner = newOwener				//转让房主
	err := db.DB.Model(&room).Update("room_owner",room.RoomOwner).Error		//数据库中保存修改
	var Err Error.Err
	if err != nil {
		Err = Error.ErrInfo(1,err,"转让错误")
	}else {
		Err =Error.ErrInfo(0,"","转让成功")
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

func updateRoom(c *gin.Context){
	var room ,newInfo model.Room
	room.RoomId = c.Query("roomId")		//获取房间号
	err := db.DB.Where("room_id = ?",room.RoomId).First(&room).Error		//查找房间
	var Err Error.Err
	if err != nil {
		Err = Error.ErrInfo(1,err,"房间未找到")
	}else {
		c.ShouldBind(&newInfo)						//获取新参数
		if newInfo.RoomName ==""{
			newInfo.RoomName = room.RoomName
		}
		if newInfo.RoomCap == 0 {
			newInfo.RoomCap = room.RoomCap
		}
		newRoom := map[string]interface{}{
			"RoomName": newInfo.RoomName,
			"RoomCap":  newInfo.RoomCap,
			"RoomAccess":newInfo.RoomAccess,
		}
		db.DB.Model(&room).Updates(newRoom)		//修改房间信息
		Err =Error.ErrInfo(0,"","修改成功")
	}
	c.Set("Error",Err)
	c.Abort()
	return
}

func dissolveRoom(c *gin.Context) {
	r, _ := c.Get("room")
	room := r.(model.Room)
	reply,_ := mq.Rdb.Get(room.RoomId).Bytes()		//获取房间内所有的用户信息
	var clients []model.Client
	json.Unmarshal(reply,&clients)
	mq.Rdb.Del(room.RoomId)					//删除房间内储存的用户
	for _,client:= range clients {
		reply,_ := mq.Rdb.Get(client.UserName).Bytes()		//得到用户下所有的房间
		var rooms []string
		json.Unmarshal(reply,&rooms)
		for k,v := range rooms{
			if v == room.RoomId{
				rooms = append(rooms[:k],rooms[k+1:]...)	//删除每个用户下的该房间信息
				data,_ := json.Marshal(rooms)
				mq.Rdb.Do("Set",client.UserName,data)
				break
			}
		}
		i, _ := model.ClientRooms.Load(client.UserName)		//更新用户下的房间信息
		Connes := i.([]*model.WebsocketConn)
		for k, v := range Connes {
			if v.RoomName == room.RoomName{
				Connes = append(Connes[:k],Connes[k+1:]...)
				model.ClientRooms.Store(v.UserName,Connes)
				v.Conn.Close()				//关闭websocket连接
				mq.Publish("system",v.RoomName,"exit"+v.UserName)		//发送退出消息关闭协程
			}
		}
	}
	db.DB.Delete(&room)			//mysql中删除房间
	Err := Error.ErrInfo(0,"","解散成功")
	c.Set("Error",Err)
	c.Abort()
	return
}

//发送消息
func sendImg(c *gin.Context) {
	i, _ := c.Get("user")
	user := i.(model.User)
	roomId := c.PostForm("roomId")		//获取房间号
	//判断用户是否在房间内
	reply, _ := mq.Rdb.Get(user.UserName).Bytes()
	var rooms []string
	json.Unmarshal(reply, &rooms)
	flag := false
	for _,v:=range rooms{
		if v == roomId{
			flag = true
			break
		}
	}
	if !flag {
		Err := Error.ErrInfo(0,"","不在房间内无法发送图片")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	room := model.Room{}
	db.DB.Where("room_id = ?",roomId).First(&room)	//查找房间

	Img, err := c.FormFile("Img")			//获取图片
	if err != nil {
		Err:=Error.ErrInfo(1,err,"获取图片错误")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	s := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100000000))	//设置随机数更改图片名
	Names := strings.Split(Img.Filename,".")
	ImgName := Names[0]+s+"."+Names[1]
	dst := fmt.Sprintf("./Imgs/%s", ImgName)
	c.SaveUploadedFile(Img, dst)			//将图片保存到本地
	file, err := os.Open(dst)				//打开图片
	defer file.Close()
	if err != nil {
		Err := Error.ErrInfo(1,err,"打开图片错误")
		c.Set("Error",Err)
		c.Abort()
		return
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("smfile", dst)		//创建表单
	io.Copy(part, file)			//将文件写入表单

	contentType := writer.FormDataContentType()		//得到Content-Type
	writer.Close()

	//上传图片到图床
	request, _ := http.NewRequest("POST", "https://sm.ms/api/v2/upload", body)
	request.Header.Add("Content-Type", contentType)
	request.Header.Add("Authorization", "6RgNFfmI8Fb5QwXjkgTm1uMf3ALsFt8l")
	client := http.Client{}
	resp, _ := client.Do(request)
	data, _ := ioutil.ReadAll(resp.Body)		//获取图床的响应
	var Resp interface{}
	json.Unmarshal(data,&Resp)
	var url string
	for k,v := range Resp.(map[string]interface{}){
		if k == "data"{
			for key,value := range v.(map[string]interface{}){
				if key == "url"{
					url = value.(string)		//得到图床返回的url链接
					break
				}
			}
		}
	}
	mq.Publish(user.UserName,room.RoomName,url)		//将链接发送到房间内的所有用户

	Err := Error.ErrInfo(0,"","发送图片成功")
	c.Set("Error",Err)
	c.Abort()
	return
}
