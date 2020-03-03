package services

import (
	"fmt"
	"gosystem/comm"
	"gosystem/dao"
	"gosystem/dataSource"
	"gosystem/models"
	"log"

	"github.com/gomodule/redigo/redis"
)

type BlackUserService interface {
	GetAll() []models.BlackUser
	CountAll() int64
	Get(id int) *models.BlackUser
	Delete(id int) error
	Update(data *models.BlackUser, columns []string) error
	Insert(data *models.BlackUser) error
	GetByUid(uid int) *models.BlackUser
	GetUserToday(uid int) []models.BlackUser
}

type blackUserService struct {
	dao *dao.BlackUserDao
}

func NewUserDayService() BlackUserService {
	return &blackUserService{
		dao: dao.NewBlackUserDao(dataSource.NewMysqlMaster()),
	}
}

func (this *blackUserService) GetAll() []models.BlackUser {
	return this.dao.GetAll()
}

func (this *blackUserService) CountAll() int64 {
	return this.dao.CountAll()
}

func (this *blackUserService) Get(id int) *models.BlackUser {
	user := this.getByCache(id)
	if user == nil || user.Id <= 0 {
		user = this.dao.Get(id)
		if user != nil && user.Id > 0 { //数据库读到了数据，再存入缓存，否则存进去无效数据
			this.setByCache(user)
		}
	}
	return user
}

func (this *blackUserService) Delete(id int) error {
	return this.dao.Delete(id)
}

func (this *blackUserService) Update(data *models.BlackUser, columns []string) error {
	// 先更新缓存,这里直接是清空该data对应的缓存数据;后面再读取会从数据库更新1个新数据到缓存
	this.updateByCache(data, columns)
	// 再更新数据
	return this.dao.Update(data, columns)
}

func (this *blackUserService) Insert(data *models.BlackUser) error {
	return this.dao.Insert(data)
}

func (this *blackUserService) GetByUid(uid int) *models.BlackUser {
	return this.dao.GetByUid(uid)
}

func (this *blackUserService) GetUserToday(uid int) []models.BlackUser {
	y, m, d := comm.NowTime().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	return this.dao.Search(uid, strDay)
}

//redis缓存
func (this *blackUserService) getByCache(id int) *models.BlackUser {
	key := fmt.Sprintf("info_suer_%v", id)
	rds := dataSource.RedisInstCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user service.getByCache HGETALL key=", key, ",error=", err)
		return nil
	}
	dataid := comm.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataid <= 0 {
		return nil
	}
	data := &models.BlackUser{
		Id:         id,
		Username:   comm.GetStringFromStringMap(dataMap, "Username", ""),
		BlackTime:  int(comm.GetInt64FromStringMap(dataMap, "BlackTime", 0)),
		RealName:   comm.GetStringFromStringMap(dataMap, "RealName", ""),
		Mobile:     comm.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIP:      comm.GetStringFromStringMap(dataMap, "SysIP", ""),
		SysStatus:  int(comm.GetInt64FromStringMap(dataMap, "SysStatus", 0)),
	}

	return data
}

//redis缓存
func (this *blackUserService) setByCache(user *models.BlackUser) {
	if user == nil || user.Id <= 0 {
		return
	}
	key := fmt.Sprintf("info_suer_%v", user.Id)
	rds := dataSource.RedisInstCache()

	params := redis.Args{key}
	params = params.Add(user.Id)
	if user.Username != "" {
		params = params.Add(user.Username)
		params = params.Add(user.SysStatus)
		params = params.Add(user.Mobile)
		params = params.Add(user.RealName)
		// ...
	}
	_, err := rds.Do("HMSET", params) //多字段批量添加
	if err != nil {
		log.Println("user service.setByCache HMSET key=", key, ",params=", params, ",error=", err)
	}
}

//redis缓存
func (this *blackUserService) updateByCache(user *models.BlackUser, column []string) {
	if user == nil || user.Id <= 0 {
		return
	}
	key := fmt.Sprintf("info_suer_%v", user.Id)
	dataSource.RedisInstCache().Do("DEL", key)
}
