package controllers

import (
	"fmt"
	"gosystem/utils"
	"log"

	"gosystem/comm"
	"gosystem/conf"
	"gosystem/models"
)

// package 私有方法
func (this *IndexController) checkUserDay(uid int, rdsnum int64) bool {
	userDayInfo := this.ServiceUserDay.GetUserToday(uid)
	if userDayInfo != nil && userDayInfo.Uid == string(uid) {
		// 今天存在抽奖记录
		if int(rdsnum) < userDayInfo.Num { //如果之前redis没有限制住，走到这里了，就需要更新redis数据库
			utils.InitUserLuckyNum(uid, int64(userDayInfo.Num))
		}

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
			Uid:        string(uid),
			DAY:        strday,
			Num:        1,
			SysCreated: comm.NowTime(),
			SysUpdated: comm.NowTime(),
		}
		err := this.ServiceUserDay.Insert(userDayInfo)
		if err != nil {
			log.Println("index_lucky_check_user_day ServiceUserDay.Insert error", err)
		}
		utils.InitUserLuckyNum(uid, 1) //创建时候，给redis初始化一个1
	}
	return true
}
