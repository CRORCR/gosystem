package models

type Gift struct {
	Id           int    `xorm:"INT pk autoincr 'id'"`
	Title        string `xorm:"VARCHAR(255) 'title'"`     //奖品名称
	PrizeNum     int    `xorm:"INT 'prize_num'"`          //奖品数量，0 无限量，>0限量，<0无奖品
	LeftNum      int    `xorm:"INT 'left_num'"`           //剩余数量
	PrizeCode    string `xorm:"VARCHAR(50) 'prize_code'"` //0-9999表示100%，0-0表示万分之一的中奖概率
	PrizeTime    int    `xorm:"DATETIME 'prize_time'"`    //发奖周期，D天
	Img          string `xorm:"VARCHAR(255) 'img'"`       //奖品图片
	DisplayOrder int    `xorm:"INT 'display_order'"`      //位置序号，小的排在前面
	Gtype        int    `xorm:"INT 'gtype'"`              //奖品类型，0 虚拟币，1 虚拟券，2 实物-小奖，3 实物-大奖
	Gdata        string `xorm:"VARCHAR(255) 'gdata'"`     //扩展数据，如：虚拟币数量
	TimeBegin    int    `xorm:"INT 'time_begin'"`         //开始时间
	TimeEnd      int    `xorm:"INT 'time_end'"`
	PrizeData    string `xorm:"MEDIUMTEXT 'prize_data'"` //发奖计划，[[时间1,数量1],[时间2,数量2]]
	PrizeBegin   int    `xorm:"DATETIME 'prize_begin'"`  //发奖计划周期的开始
	PrizeEnd     int    `xorm:"DATETIME 'prize_end'"`
	SysCreated   int    `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP created 'sys_created'"`
	SysUpdated   int    `xorm:"DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP updated 'sys_updated'"`
	SysStatus    int    `xorm:"SMALLINT 'sys_status'"` //状态，0 正常，1 删除'
	SysIP        string `xorm:"VARCHAR(50) 'sys_ip'"`  //操作人IP
}

func (this *Gift) TableName() string {
	return "gift"
}
