package models

type GiftPrize struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	PrizeNum     int    `json:"-"` //奖品数量：0 无限；>0 限量；<0 无奖品
	LeftNum      int    `json:"-"`
	PrizeCodeA   int    `json:"-"`
	PrizeCodeB   int    `json:"-"`
	Img          string `json:"img"`
	DisplayOrder int    `json:"display_order"`
	Gtype        int    `json:"gtype"`
	Gdata        string `json:"gdata"`
}
