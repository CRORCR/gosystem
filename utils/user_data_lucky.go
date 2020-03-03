package utils

import (
	"fmt"
	"gosystem/comm"
	"gosystem/dataSource"
	"log"
	"math"
	"time"
)

const UserFrameSize = 2

func init() {
	restGroupUserList()

}

func restGroupUserList() {
	log.Println("user_data_lucky.restGroupUserList start")
	rds := dataSource.RedisInstCache()
	for i := 0; i < UserFrameSize; i++ {
		key := fmt.Sprintf("day_users_%v", i)
		rds.Do("DEL", key)
	}
	log.Println("user_data_lucky.restGroupUserList stop")
	//零点定时任务
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, restGroupUserList)
}

func IncrUserLuckyNum(uid int) int64 {
	i := uid % UserFrameSize
	key := fmt.Sprintf("day_users_%v", i)
	rds := dataSource.RedisInstCache()
	rs, err := rds.Do("HINCRBY", key, uid, 1) //hincrby 将名称为key的hash中field的value增加integer
	if err != nil {
		log.Printf("user_data_lucky redis IncrUserLuckyNum hincrby key[%v] uid[%v] err[%v]", key, uid, err)
		return math.MaxInt32
	}
	return rs.(int64)
}

//每日参与次数要以数据库为准
// 从数据库初始化缓存数据
func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % UserFrameSize
	key := fmt.Sprintf("day_users_%v", i)
	rds := dataSource.RedisInstCache()
	_, err := rds.Do("HSET", key, uid, num)
	if err != nil {
		log.Printf("user_data_lucky redis InitUserLuckyNum hset key[%v] uid[%v] num[%v] err[%v]", key, uid, num, err)
		return
	}

}
