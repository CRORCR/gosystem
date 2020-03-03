package services

import (
	"fmt"
	"gosystem/comm"
	"gosystem/conf"
	"gosystem/dao"
	"gosystem/dataSource"
	"gosystem/models"
	"log"

	"github.com/gomodule/redigo/redis"
)

type BlackIpService interface {
	GetAll() []models.BlackIp
	CountAll() int64
	Get(id int) *models.BlackIp
	Delete(id int) error
	Update(data *models.BlackIp, columns []string) error
	Insert(data *models.BlackIp) error
	GetByIp(ip string) *models.BlackIp
}

type blackIpService struct {
	dao *dao.BlackIpDao
}

func NewBlackIpService() BlackIpService {
	return &blackIpService{
		dao: dao.NewBlackIpDao(dataSource.NewMysqlMaster()),
	}
}

func (this *blackIpService) GetAll() []models.BlackIp {
	return this.dao.GetAll()
}

func (this *blackIpService) CountAll() int64 {
	return this.dao.CountAll()
}

func (this *blackIpService) Search(ip string) *models.BlackIp {
	return this.dao.GetByIp(ip)
}

func (this *blackIpService) Get(id int) *models.BlackIp {
	return this.dao.Get(id)
}

func (this *blackIpService) Delete(id int) error {
	return this.dao.Delete(id)
}

func (this *blackIpService) Update(data *models.BlackIp, columns []string) error {
	return this.dao.Update(data, columns)
}

func (this *blackIpService) Insert(data *models.BlackIp) error {
	return this.dao.Insert(data)
}

func (s *blackIpService) GetByIp(ip string) *models.BlackIp {
	data := s.getByCache(ip)
	if data == nil || data.Ip == "" {
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.BlackIp{Ip: ip}
		}
		s.setByCache(data)
	}
	return data
}

func (s *blackIpService) getByCache(ip string) *models.BlackIp {
	key := fmt.Sprintf(conf.RdsBlackipCacheKeyPrefix+"%s", ip)
	rds := dataSource.RedisInstCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("blackip_service.getByCache HGETALL key = ", key, ", error = ", err)
		return nil
	}
	dataIp := comm.GetStringFromStringMap(dataMap, "Ip", "")
	if dataIp == "" {
		return nil
	}
	data := &models.BlackIp{
		Id:         int(comm.GetInt64FromStringMap(dataMap, "Id", 0)),
		Ip:         dataIp,
		BlackTime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
	}
	return data
}

func (s *blackIpService) setByCache(data *models.BlackIp) {
	if data == nil || data.Ip == "" {
		return
	}

	key := fmt.Sprintf(conf.RdsBlackipCacheKeyPrefix+"%s", data.Ip)
	rds := dataSource.RedisInstCache()
	// 数据更新到redis缓存
	params := []interface{}{key}
	params = append(params, "Ip", data.Ip)
	if data.Id > 0 {
		params = append(params, "Blacktime", data.BlackTime)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
	}
	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("blackip_service.setByCache HMSET params = ", params, ", error = ", err)
	}
}

func (s *blackIpService) updateByCache(data *models.BlackIp, columns []string) {
	if data == nil || data.Ip == "" {
		return
	}
	key := fmt.Sprintf(conf.RdsBlackipCacheKeyPrefix+"%s", data.Ip)
	rds := dataSource.RedisInstCache()
	rds.Do("DEL", key)
}
