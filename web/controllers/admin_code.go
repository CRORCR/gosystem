package controllers

import (
	"github.com/kataras/iris"

	"gosystem/services"
)

type AdminCodeController struct {
	Ctx         iris.Context
	ServiceCode services.CodeService
}
