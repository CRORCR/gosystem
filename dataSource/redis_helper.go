package dataSource

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"gosystem/conf"
)

var redisInst *RedisConn
var redisLock sync.Mutex

type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

func (this *RedisConn) Do(commandName string,
	args ...interface{}) (reply interface{}, err error) {

	conn := this.pool.Get() //连接池拿一个连接
	defer conn.Close()      // 将连接放回连接池

	t1 := time.Now().UnixNano()
	reply, err = conn.Do(commandName, args...)

	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Fatal("redis_helper.Do error ", err, e)
		}
	}

	t2 := time.Now().UnixNano()

	//是否需要打印连接日志
	if this.showDebug {
		fmt.Printf(
			"[redis] [info] [%dus] cmd=%s, args=%v, reply=%s, err=%s\n",
			(t2-t1)/1000, commandName, args, reply, err,
		)
	}

	return reply, err

}

func (this *RedisConn) ShowDebug(show bool) {
	this.showDebug = show
}

//单利模式
func RedisInstCache() *RedisConn {

	if redisInst != nil {
		return redisInst
	}

	redisLock.Lock()
	defer redisLock.Unlock()

	//两次判断 小细节要注意
	if redisInst != nil {
		return redisInst
	}

	return NewRedisCache()

}

func NewRedisCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp",
				fmt.Sprintf("%s:%d",
					conf.RedisCache.Host,
					conf.RedisCache.Port,
				),
			)

			if err != nil {
				log.Fatal("redis_helper.NewRedisCache error ", err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000, // 最多连接数
		MaxActive:       10000, // 最多活跃数
		IdleTimeout:     0,     // 超时时间
		Wait:            false, // 连接等待
		MaxConnLifetime: 0,     //最大连接时间，0 一直连接
	}

	redisInst = &RedisConn{pool: &pool}
	redisInst.ShowDebug(true)
	return redisInst
}
