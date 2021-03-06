package utils

import (
	"fmt"
	"log"
	"math"
	"time"

	"gosystem/comm"
	"gosystem/dataSource"
)

const ipFrameSize = 2

// init 程序启动时调用1次 清理ip分段数据
func init() {
	resetGroupIpList()
}

//每日凌晨数据清空
func resetGroupIpList() {
	log.Println("ip_day_lucky.resetGroupIpList start")
	redisDB := dataSource.RedisInstCache()
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_ips_%d", i)
		_, _ = redisDB.Do("DEL", key)
	}
	log.Println("ip_day_lucky.resetGroupIpList stop")

	// IP 当天的统计数，零点的时候归零，设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIpList)
}

//今天的IP抽奖次数递增，返回递增后的数值 原子性递增
func IncrIpLuckyNum(strIp string) int64 {
	ip := comm.Ip4ToInt(strIp)
	i := ip % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i) //散列存储
	redisDB := dataSource.RedisInstCache()
	rs, err := redisDB.Do("HINCRBY", key, ip, 1) //hash结构存储
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error ", err)
		return math.MaxInt64
	}
	return rs.(int64)
}
