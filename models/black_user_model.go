package models

type BlackUser struct {
	Id         int    `xorm:"INT pk autoincr 'id'"`
	Uid        int    `xorm:"INT 'uid'"`
	Username   string `xorm:"VARCHAR(50) 'username'"`
	BlackTime  int    `xorm:"INT 'black_time'"`
	RealName   string `xorm:"VARCHAR(50) 'real_name'"`
	Mobile     string `xorm:"VARCHAR(50) 'mobile'"`
	Address    string `xorm:"VARCHAR(255) 'address'"`
	SysCreated int    `xorm:"INT  'sys_created'"`
	SysUpdated int    `xorm:"INT  'sys_updated'"`
	SysIP      string `xorm:"VARCHAR(50) 'sys_ip'"`
	SysStatus  int    `xorm:"SMALLINT 'sys_status'"`
}

func (this *BlackUser) TableName() string {
	return "black_user"
}
