package controllers

import (
	"fmt"
	"log"

	"gosystem/comm"
	"gosystem/conf"
	"gosystem/models"
)

func (this *IndexController) checkUserDay(uid int) bool {
	userDayInfo := this.ServiceUserDay.GetUserToday(uid)
	if userDayInfo != nil && userDayInfo.Uid == uid {
		// 今天存在抽奖记录
		if userDayInfo.Num >= conf.UserPrizeMax {
			return false
		} else {
			userDayInfo.Num++
			userDayInfo.SysUpdated = comm.NowTime()
			err := this.ServiceUserDay.Update(
				userDayInfo,
				[]string{
					"num",
					"sys_updated",
				},
			)
			if err != nil {
				log.Println("index_lucky_check_user_day ServiceUserDay.Update error", err)
			}
		}
	} else {
		// 创建今天的用户参与记录
		y, m, d := comm.NowTime().Date()             //获得年月日
		strday := fmt.Sprintf("%d%02d%02d", y, m, d) //转为字符串
		userDayInfo = &models.UserDay{
			Uid:        uid,
			DAY:        strday,
			Num:        1,
			SysCreated: comm.NowTime(),
			SysUpdated: comm.NowTime(),
		}
		err := this.ServiceUserDay.Insert(userDayInfo)
		if err != nil {
			log.Println("index_lucky_check_user_day ServiceUserDay.Insert error", err)
		}
	}
	return true
}
