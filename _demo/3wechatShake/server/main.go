/**
 * 微信摇一摇
 * 基础功能
 * /lucky 只有一个抽奖的接口
 * 压力测试
 * wrk -t10 -c10 -d5 http://localhost:8080/lucky
	手机改成了两万，中奖概率100%。那么如果压测，应该是两万条记录，奖品都是手机

	清空数据 echo "" > /lottery/_demo/3wechatShake/log/lottery_demo.log
	wc -l lottery/_demo/3wechatShake/log/lottery_demo.log 查看日志多少行日志

	测试结果显示超发了，这是很严重的问题，线程不安全


 */
package server

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type lotteryController struct {
	Ctx iris.Context
}

// 返回数量的信息 GET http://localhost:8080/
func (c *lotteryController) Get() string {
	count := 0
	total := 0
	for _, data := range giftList {
		if data.isUse && (data.total == 0 || (data.total > 0 && data.left > 0)) {
			count++
			total += data.left
		}
	}

	return fmt.Sprintf("当前有效奖品种类数量：%d，限量奖品总数量：%d\n", count, total)
}

// 返回数量的信息 GET http://localhost:8080//lucky
func (c *lotteryController) GetLucky() map[string]interface{} {
	//锁是为了结果数据共享，只有共享的数据才需要加锁 在这里就是奖品切片 giftList
	mu.Lock()
	defer mu.Unlock()

	code := luckyCode()

	result := map[string]interface{}{}

	ok := false
	result["success"] = false

	for _, data := range giftList {
		if !data.isUse || (data.total > 0 && data.left <= 0) {
			continue
		}

		// 中奖了，抽奖编码在奖品编码范围内
		if data.rateMin <= int(code) && int(code) < data.rateMax {

			// 开始发奖
			sendData := ""
			switch data.gtype {
			case giftTypeCoin:
				ok, sendData = sendCoin(data)
			case giftTypeCoupon:
				ok, sendData = sendCoupon(data)
			case giftTypeCouponFix:
				ok, sendData = sendCouponFix(data)
			case giftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			case giftTypeRealLarge:
				ok, sendData = sendRealLarge(data)
			}

			if ok {
				// 中奖成功，成功得到奖品
				// 生成中奖记录
				saveLuckyData(code, sendData, data)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				break
			}

		}
	}

	if v, ok := result["success"]; ok && v == false {
		result["data"] = "没有中奖"
	}

	return result
}

func NewApp() *iris.Application {
	app := iris.Default()
	mvc.New(app.Party("/")).Handle(&lotteryController{})

	initLog()
	initGift()

	return app
}

func Run() {
	app := NewApp()
	app.Run(iris.Addr(":8080"))
}
