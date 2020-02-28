package controllers

import (
	"gosystem/conf"
	"gosystem/models"
)

func (this *IndexController) prize(prizeCode int, limitBlack bool) *models.GiftPrize {
	var prizeGift *models.GiftPrize

	giftList := this.ServiceGift.GetAllUse()

	for _, gift := range giftList {

		// 中奖编码区间满足条件，说明可以中奖
		if gift.PrizeCodeA <= prizeCode && gift.PrizeCodeB >= prizeCode {
			//这个判断条件用了或的关系，跟我的风格不一样
			//是黑名单 检查是不是小奖 是 发放
			//不是黑名单 直接发
			if !limitBlack || gift.Gtype < conf.GiftTypeGiftSmall {
				prizeGift = &gift
				break
			}
		}
	}

	return prizeGift
}
