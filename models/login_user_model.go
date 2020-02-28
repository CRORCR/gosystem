package models

import "time"

type LoginUser struct {
	Uid      int
	Username string
	Now      time.Time
	Ip       string
	Sign     string //签名数据
}

//基于cookie的用户状态
//ObjLoginUser登陆用户对象
//登陆用户与cookie读写
//cookie的安全校验值，不能被篡改
