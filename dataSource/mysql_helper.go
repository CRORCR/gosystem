package dataSource

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"gosystem/conf"
)

var MysqlMasterInst *xorm.Engine
var mysqlLock sync.Mutex

//单例模式 - 得到唯一的主库实例
//  在应用运行期间会不断调用数据库操作,不能每次调用都实例化1次
func MysqlInstMaster() *xorm.Engine {
	if MysqlMasterInst != nil {
		return MysqlMasterInst
	}

	// 处理高并发时避免重复定义实例
	mysqlLock.Lock()
	defer mysqlLock.Unlock()

	// 锁定后 return NewDbMaster() 直接创建可能也会出问题
	// 有1个在创建,后面2个排队;这个创建完了后面2个再进来,又实例化,也不行
	// 导致失败单例

	// 还要再判断是否创建
	if MysqlMasterInst != nil {
		return MysqlMasterInst
	}

	return NewMysqlMaster()
}

// 返回xorm的MySQL数据库操作引擎
func NewMysqlMaster() *xorm.Engine {
	engine, err := xorm.NewEngine(
		conf.MysqlDriverName,
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
			conf.DbMaster.User,
			conf.DbMaster.Pwd,
			conf.DbMaster.Host,
			conf.DbMaster.Port,
			conf.DbMaster.Database,
		),
	)

	if err != nil {
		log.Fatal("db_helper.NewMysqlMaster NewEngine error ", err)
		return nil
	}

	// xorm支持的调试特性
	// SQL执行时间
	// instance.ShowExecTime()
	// 执行的SQL语句 生产false不展示 开发true展示
	//instance.ShowSQL(false)
	// 本地调试打开 SQL 调试
	engine.ShowSQL(true)

	MysqlMasterInst = engine

	return MysqlMasterInst
}
