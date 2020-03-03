package services

import (
	"fmt"
	"gosystem/dataSource"
	"strconv"
	"time"

	"gosystem/dao"
	"gosystem/models"
)

type UserdayService interface {
	GetAll(page, size int) []models.UserDay
	CountAll() int64
	Search(uid int, day int) []models.UserDay
	Count() int64
	Get(id int) *models.UserDay
	Update(user *models.UserDay, columns []string) error
	Create(user *models.UserDay) error
	GetUserToday(uid int) *models.UserDay
}

type userdayService struct {
	dao *dao.UserDayDao
}

func NewUserdayService() UserdayService {
	return &userdayService{
		dao: dao.NewUserDayDao(dataSource.MysqlInstMaster()),
	}
}

func (s *userdayService) GetAll(page, size int) []models.UserDay {
	return s.dao.GetAll(page, size)
}

func (s *userdayService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userdayService) Search(uid int, day int) []models.UserDay {
	return s.dao.Search(uid, day)
}

func (s *userdayService) Count() int64 {
	return s.dao.CountAll()
}

func (s *userdayService) Get(id int) *models.UserDay {
	return s.dao.Get(id)
}

func (s *userdayService) Update(data *models.UserDay, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *userdayService) Create(data *models.UserDay) error {
	return s.dao.Insert(data)
}

func (s *userdayService) GetUserToday(uid int) *models.UserDay {
	y, m, d := time.Now().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	day, _ := strconv.Atoi(strDay)
	list := s.dao.Search(uid, day)
	if list != nil && len(list) > 0 {
		return &list[0]
	} else {
		return nil
	}
}
