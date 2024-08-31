package mysqltool

import (
	"context"
	"github.com/yunduansing/gtools/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order struct {
	OrderId string       `gorm:"primarykey;type:varchar(50);" json:"orderId"`
	Amount  int64        `json:"amount"`
	Goods   []OrderGoods `json:"goods"`
}

type OrderGoods struct {
	OrderId string `gorm:"primarykey;type:varchar(50);" json:"orderId"`
	SkuId   string `gorm:"primarykey;type:varchar(50);" json:"skuId"`
	Num     int    `json:"num"`
	Price   int64  `json:"price"`
}

func initOrders(db *gorm.DB) {
	orders := []Order{
		{
			OrderId: "1",
			Amount:  250,
			Goods:   []OrderGoods{},
		},
		{
			OrderId: "2",
			Amount:  300,
			Goods:   []OrderGoods{},
		},
	}

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: clause.PrimaryKey}},
		DoNothing: true,
	}).Create(&orders)

	goods := []OrderGoods{
		{
			OrderId: "2",
			SkuId:   "3",
			Num:     1,
			Price:   200,
		},
		{
			OrderId: "1",
			SkuId:   "1",
			Num:     1,
			Price:   100,
		},
		{
			OrderId: "1",
			SkuId:   "2",
			Num:     1,
			Price:   150,
		},
	}

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: clause.PrimaryKey}},
		DoNothing: true,
	}).Create(&goods)
}

func findAllOrders(db *gorm.DB) {
	var orders []Order
	err := db.Preload("Goods").Find(&orders).Error
	if err != nil {
		logger.GetLogger().Error(context.Background(), "find orders err:", err)
		return
	}
	logger.GetLogger().Info(context.Background(), "find orders list:", orders)
}
