package models

import "time"

type Result struct {
	Id         int       `xorm:"INT pk autoincr 'id'"`
	GiftId     int       `xorm:"INT 'gift_id'"`            //奖品ID，关联lt_gift表
	GiftName   string    `xorm:"VARCHAR(255) 'gift_name'"` //奖品名称
	GiftType   int       `xorm:"INT 'gift_type'"`          //奖品类型，同Gift. gtype
	Uid        int       `xorm:"INT 'uid'"`                //用户ID
	Username   string    `xorm:"VARCHAR(50) 'username'"`   //用户名
	PrizeCode  int       `xorm:"INT 'prize_code'"`         //抽奖编号（4位的随机数）
	GiftData   string    `xorm:"VARCHAR(50) 'gift_data'"`  //获奖信息
	SysCreated time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP created 'sys_created'"`
	SysStatus  int       `xorm:"SMALLINT 'sys_status'"` //状态，0 正常，1删除，2作弊
	SysIP      string    `xorm:"VARCHAR(50) 'sys_ip'"`  //用户抽奖的IP
}

func (this *Result) TableName() string {
	return "result"
}
