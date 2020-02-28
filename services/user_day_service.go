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

type UserDayService interface {
	GetAll() []models.UserDay
	CountAll() int64
	Get(id int) *models.BlackUser
	Delete(id int) error
	Update(data *models.UserDay, columns []string) error
	Insert(data *models.UserDay) error
	GetByUid(uid int) *models.UserDay
	GetUserToday(uid int) *models.UserDay
}

type userDayService struct {
	dao *dao.UserDayDao
}

func NewUserDayService() UserDayService {
	return &userDayService{
		dao: dao.NewUserDayDao(dataSource.NewMysqlMaster()),
	}
}

func (this *userDayService) GetAll() []models.UserDay {
	return this.dao.GetAll()
}

func (this *userDayService) CountAll() int64 {
	return this.dao.CountAll()
}

func (this *userDayService) Get(id int) *models.BlackUser {
	user := this.getByCache(id)
	if user == nil || user.Id <= 0 {
		user = this.dao.Get(id)
		if user != nil && user.Id > 0 { //数据库读到了数据，再存入缓存，否则存进去无效数据
			this.setByCache(user)
		}
	}
	return user
}

func (this *userDayService) Delete(id int) error {
	return this.dao.Delete(id)
}

func (this *userDayService) Update(data *models.UserDay, columns []string) error {
	return this.dao.Update(data, columns)
}

func (this *userDayService) Insert(data *models.UserDay) error {
	return this.dao.Insert(data)
}

func (this *userDayService) GetByUid(uid int) *models.UserDay {
	return this.dao.GetByUid(uid)
}

func (this *userDayService) GetUserToday(uid int) *models.UserDay {
	y, m, d := comm.NowTime().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	return this.dao.Search(uid, strDay)
}

//redis缓存
func (this *userDayService) getByCache(id int) *models.BlackUser {
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
func (this *userDayService) setByCache(user *models.BlackUser) {
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
func (this *userDayService) updateByCache(user *models.BlackUser, column []string) {
	if user == nil || user.Id <= 0 {
		return
	}
	key := fmt.Sprintf("info_suer_%v", user.Id)
	dataSource.RedisInstCache().Do("DEL", key)
}
