package controllers

import (
	"fmt"
	"gosystem/comm"
	"gosystem/conf"
	"gosystem/models"
	"gosystem/utils"
	"strings"

	"github.com/kataras/iris"

	"gosystem/services"
)

type AdminCodeController struct {
	Ctx         iris.Context
	ServiceCode services.CodeService
}

func (c *AdminCodeController) PostImport() {
	giftId, _ := c.Ctx.URLParamInt("gift_id")
	if giftId < 1 {
		c.Ctx.Text("没有指定奖品ID，无法进行导入,<a href='' onclick='history.go(-1)'")
		return
	}
	gift := c.ServiceCode.Get(giftId, false) //是否从redis获取
	if gift == nil || gift.Id < 1 || gift.Gtype != conf.GiftTypeCodeDiff {
		c.Ctx.Text("奖品信息不存在，或者奖品类型不是差异化优惠券，无法进行导入")
		return
	}

	codes := c.Ctx.PostValue("codes")
	now := comm.NowUnix()
	list := strings.Split(codes, "\n")
	sucNum := 0
	errNum := 0

	for _, code := range list {
		code := strings.TrimSpace(code)
		if code == "" {
			continue
		}
		data := &models.Code{
			GiftId:     giftId,
			Code:       code,
			SysCreated: now,
		}
		err := c.ServiceCode.Created(data)
		if err != nil {
			errNum++
		} else {
			sucNum++
			//成功导入数据库后，再存入缓存一份
			utils.ImportCacheCodes(giftId, code)
		}
	}
}

func (c *AdminCodeController) Getrecache() {
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	id, err := c.Ctx.URLParamInt("id")
	if id < 1 || err != nil {
		rs := fmt.Sprintf("没有指定优惠券所属的奖品ID")
		c.Ctx.HTML(rs)
		return
	}
	sucNum, errNum := utils.RecacheCodes(id, c.ServiceCode)
	rs := fmt.Sprintf("sucNum=%v errNum=%v <a href='%s'></a>", sucNum, errNum, refer)
	c.Ctx.HTML(rs)
}
