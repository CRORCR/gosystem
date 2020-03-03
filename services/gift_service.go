package services

import (
	"encoding/json"
	"gosystem/comm"
	"gosystem/dao"
	"gosystem/dataSource"
	"gosystem/models"
	"log"
	"strconv"
	"strings"
)

type GiftService interface {
	GetAll(useCache bool) []models.Gift
	CountAll() int64
	Get(id int) *models.Gift
	Delete(id int) error
	Update(data *models.Gift, columns []string) error
	Insert(data *models.Gift) error
	GetAllUse() []models.GiftPrize
	DecrLeftNum(id, num int) (int64, error)
	IncrLeftNum(id, num int) (int64, error)
}

type giftService struct {
	dao *dao.GiftDao
}

// 返回 GiftService 接口 而不是 私有的giftService 否则外界无法使用
func NewGiftService() GiftService {
	return &giftService{
		dao: dao.NewGiftDao(dataSource.NewMysqlMaster()),
	}
}

//参数判断是否需要缓存
func (this *giftService) GetAll(useCache bool) []models.Gift {
	if !useCache {
		return this.dao.GetAll()
	}
	//如果使用缓存
	cache := this.getAllByCache()
	if len(cache) == 0 {
		cache = this.dao.GetAll()
		this.setAllByCache(cache)
	}
	return cache

}

func (this *giftService) CountAll() int64 {
	gifts := this.GetAll(true)
	return int64(len(gifts))
}

func (this *giftService) Get(id int) *models.Gift {
	gifts := this.GetAll(true)
	for _, value := range gifts {
		if value.Id == id {
			return &value
		}
	}
	return nil
}

//删除和更新都要更新缓存数据 保证一致性
func (this *giftService) Delete(id int) error {
	gift := &models.Gift{Id: id}
	this.updateByCache(gift, nil)
	return this.dao.Delete(id)
}

//删除和更新都要更新缓存数据 保证一致性
func (this *giftService) Update(data *models.Gift, columns []string) error {
	this.updateByCache(data, columns)
	return this.dao.Update(data, columns)
}

//删除和更新都要更新缓存数据 保证一致性
func (this *giftService) Insert(data *models.Gift) error {
	this.updateByCache(data, nil)
	return this.dao.Insert(data)
}

//获得所有有效的奖品 时间满足，状态正常
//gtype倒序 display_order升序
func (this *giftService) GetAllUse() []models.GiftPrize {
	dataList := make([]models.Gift, 0)

	now := comm.NowUnix()
	gifts := this.GetAll(true)
	for _, gift := range gifts {
		if gift.Id > 0 && gift.SysStatus == 0 && gift.PrizeNum >= 0 &&
			gift.TimeBegin >= now && gift.TimeEnd <= now {
			dataList = append(dataList, gift)
		}
	}

	if len(dataList) == 0 {
		return []models.GiftPrize{}
	}
	giftList := make([]models.GiftPrize, 0)

	for _, gift := range dataList {
		codes := strings.Split(gift.PrizeCode, "-")
		if len(codes) != 2 {
			continue
		}
		a, e1 := strconv.Atoi(codes[0])
		b, e2 := strconv.Atoi(codes[1])
		if e1 == nil && e2 == nil && b >= a && a >= 0 && b <= 10000 {
			data := models.GiftPrize{
				Id:           gift.Id,
				Title:        gift.Title,
				PrizeNum:     gift.PrizeNum,
				LeftNum:      gift.LeftNum,
				PrizeCodeA:   a,
				PrizeCodeB:   b,
				Img:          gift.Img,
				DisplayOrder: gift.DisplayOrder,
				Gtype:        gift.Gtype,
				Gdata:        gift.Gdata,
			}
			giftList = append(giftList, data)
		}
	}
	return giftList
}

func (this *giftService) DecrLeftNum(id, num int) (int64, error) {
	return this.dao.DecrLeftNum(id, num)
}

func (this *giftService) IncrLeftNum(id, num int) (int64, error) {
	return this.dao.IncrLeftNum(id, num)
}

//redis缓存
func (this *giftService) getAllByCache() []models.Gift {
	key := "allgift"
	rds := dataSource.RedisInstCache()
	rs, err := rds.Do("GET", key)
	if err != nil {
		log.Println("gift service.getAllByCache get key=", key, ",error=", err)
		return nil
	}
	s := comm.GetString(rs, "")
	if s == "" {
		return nil
	}
	dataList := []map[string]interface{}{}
	json.Unmarshal([]byte(s), &dataList)
	gifts := make([]models.Gift, len(dataList))
	for i := 0; i < len(dataList); i++ {
		data := dataList[i]
		id := comm.GetInt64FromMap(data, "Id", 0)
		if id <= 0 {
			gifts[i] = models.Gift{}
			continue
		}
		gift := models.Gift{
			Id:           int(id),
			Title:        comm.GetStringFromMap(data, "Title", ""),
			PrizeNum:     int(comm.GetInt64FromMap(data, "prize_num", 0)),
			LeftNum:      0,
			PrizeCode:    "",
			PrizeTime:    0,
			Img:          "",
			DisplayOrder: 0,
			Gtype:        0,
			Gdata:        "",
			TimeBegin:    0,
			TimeEnd:      0,
			PrizeData:    "", //发奖计划，这里不要序列化存储，因为这个字段非常大
			PrizeBegin:   0,
			PrizeEnd:     0,
			SysCreated:   0,
			SysStatus:    0,
			SysIP:        "",
		}
		gifts = append(gifts, gift)
	}
	return gifts
}

// 将奖品的数据更新到Redis缓存
func (this *giftService) setAllByCache(gifts []models.Gift) {
	strValue := ""
	if len(gifts) > 0 {
		dataList := make([]map[string]interface{}, len(gifts))
		//数据结构 []models.Gift 转换为 []map[string]interface{}{}
		for i := 0; i < len(gifts); i++ {
			gift := gifts[i]
			data := make(map[string]interface{})
			data["Id"] = gift.Id
			data["Title"] = gift.Title
			//....
			dataList = append(dataList, data)
		}
		str, err := json.Marshal(dataList)
		if err != nil {
			log.Println("gift_service.setAllByCache json.Marshal error=", err)
			return
		}
		strValue = string(str)
		key := "allgift"
		_, err = dataSource.RedisInstCache().Do("SET", key, strValue)
		if err != nil {
			log.Println("gift_service.setAllByCache redis set key=", key, ",error=", err)
		}
	}
}

//redis缓存
func (this *giftService) updateByCache(data *models.Gift, column []string) {
	if data == nil || data.Id <= 0 {
		return
	}
	key := "allgift"
	dataSource.RedisInstCache().Do("DEL", key)
}
