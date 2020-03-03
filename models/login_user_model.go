package models

type LoginUser struct {
	Uid      int
	Username string
	Now      int // 时间戳
	Ip       string
	Sign     string // 签名,签名生成 验证 cookie识别 序列化保存
}

//基于cookie的用户状态
//ObjLoginUser登陆用户对象
//登陆用户与cookie读写
//cookie的安全校验值，不能被篡改
