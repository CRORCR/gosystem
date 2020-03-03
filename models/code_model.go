package models

import "time"

type Code struct {
	Id         int       `xorm:"INT pk autoincr 'id'"`
	GiftId     int       `xorm:"INT  'gift_id'"`      //奖品 ID，关联 gift 表
	Code       string    `xorm:"VARCHAR(255) 'code'"` //虚拟券编码
	SysCreated time.Time `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP created 'sys_created'"`
	SysUpdated int       `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP updated 'sys_updated'"`
	SysStatus  int       `xorm:"SMALLINT 'sys_status'"` //状态：0 正常；1 作废；2 已发放
}

func (this *Code) TableName() string {
	return "code"
}
