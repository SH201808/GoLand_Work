package User

//请求req
type GetUserReq struct {
	UserID int64 `json:"userId"`
}

//返回resp
type GetUserResp struct{
	UserID int64 `json:"userId"`
	UserName string `json:"userName"`
}

