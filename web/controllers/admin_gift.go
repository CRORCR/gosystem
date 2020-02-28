package controllers

import (
	"github.com/kataras/iris"

	"gosystem/services"
)

type AdminGiftController struct {
	Ctx         iris.Context
	ServiceGift services.GiftService
}
