package models

type BlackUser struct {
	Id         int    `xorm:"INT pk autoincr 'id'"`
	Uid        int    `xorm:"INT 'uid'"`
	Username   string `xorm:"VARCHAR(50) 'username'"`  //用户名
	BlackTime  int    `xorm:"INT 'black_time'"`        //黑名单限制到期时间
	RealName   string `xorm:"VARCHAR(50) 'real_name'"` //联系人
	Mobile     string `xorm:"VARCHAR(50) 'mobile'"`    //手机号
	Address    string `xorm:"VARCHAR(255) 'address'"`  //联系地址
	SysCreated int    `xorm:"INT  'sys_created'"`
	SysUpdated int    `xorm:"INT  'sys_updated'"`
	SysIP      string `xorm:"VARCHAR(50) 'sys_ip'"` //IP地址
	SysStatus  int    `xorm:"SMALLINT 'sys_status'"`
}

func (this *BlackUser) TableName() string {
	return "black_user"
}
