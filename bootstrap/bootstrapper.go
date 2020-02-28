package bootstrap

import (
	"log"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"gosystem/comm"
	"gosystem/conf"
)

const (
	StaticAssets = "./public/"   //站点对外目录
	Favicon      = "favicon.ico" //sql文件
)

type Configurator func(bootstrapper *Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnData time.Time
}

func New(appName, appOwner string, cfgList ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		Application:  iris.New(),
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnData: comm.NowTime(),
	}

	for _, cfg := range cfgList {
		cfg(b)
	}

	return b
}

//初始化
func (this *Bootstrapper) Bootstrap() *Bootstrapper {
	//this.SetupViews("./views") //设置模版
	this.SetupErrorHandler() //设置异常信息

	this.Favicon(StaticAssets + Favicon) //设置默认图标
	//this.StaticWeb(StaticAssets[1:len(StaticAssets)-1], StaticAssets) //设置静态站点

	this.setupCron()

	this.Use(recover.New())
	this.Use(logger.New())

	return this
}

//监听
func (this *Bootstrapper) Listen(addr string, cfgList ...iris.Configurator) {
	err := this.Run(iris.Addr(addr), cfgList...)

	if err != nil {
		log.Fatal("bootstrap.Listen error ", err)
	}
}

//模版初始化
func (this *Bootstrapper) SetupViews(viewDir string) {
	htmlEngine := iris.HTML(viewDir, ".html").Layout("shared/layout.html")
	//htmlEngine := iris.HTML(viewDir, ".html")

	//测试环境每次修改模版都会加载，修改比较方便，生产环境记得设置 false
	htmlEngine.Reload(true)

	htmlEngine.AddFunc("FromUnixTimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeFormShort)
	})

	htmlEngine.AddFunc("FromUnixTime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeForm)
	})

	this.RegisterView(htmlEngine)
}

//异常处理
func (this *Bootstrapper) SetupErrorHandler() {
	this.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     this.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		//如果没有json输出，就用模版输出
		//ctx.ViewData("Err", err)
		//ctx.ViewData("Title", "Error")
		//ctx.View("shared/error.html")
	})
}

func (this *Bootstrapper) Configure(cfgList ...Configurator) {
	for _, cfg := range cfgList {
		cfg(this)
	}
}

//计划任务
func (this *Bootstrapper) setupCron() {
	// TODO
}
