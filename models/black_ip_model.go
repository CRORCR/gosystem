package models

type BlackIp struct {
	Id         int    `xorm:"INT pk autoincr 'id'"`
	Ip         string `xorm:"VARCHAR(5) 'ip'"`       //IP地址
	BlackTime  int    `xorm:"DATETIME 'black_time'"` //黑名单限制到期时间
	SysCreated int    `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP created 'sys_created'"`
	SysUpdated int    `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP updated 'sys_updated'"`
	SysStatus  int    `xorm:"SMALLINT 'sys_status'"`
}

func (this *BlackIp) TableName() string {
	return "black_ip"
}
