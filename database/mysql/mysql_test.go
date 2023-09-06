package mysqltool

import (
	"fmt"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestNewMySQLFromConfig(t *testing.T) {
	db, err := NewMySQLFromConfig(&Config{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		DbName:   "shop_order",
		MaxConn:  0,
		IdleConn: 0,
		LogFile:  "log/db",
	})
	if err != nil {
		logger.Error("mysql init err", err)
		return
	}

	db.AutoMigrate(&Order{})

	queryOrder(db)
	orderUpdateState(db)
}

type Order struct {
	OrderId        string `gorm:"primarykey;type:varchar(50);"`
	UserId         int64  `gorm:"type:bigint;"`
	Amount         int64  `gorm:"type:bigint;"`
	Freight        int64  `gorm:"type:bigint;"`
	DiscountAmount int64  `gorm:"type:bigint;"`
	State          int    `gorm:"type:tinyint(4);"`
	Msg            string `gorm:"type:varchar(1000)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func queryOrder(db *gorm.DB) {
	var orders []Order
	db.FindInBatches(&orders, 10, func(tx *gorm.DB, batch int) error {

		logger.Info("queryOrder", orders, batch)
		return nil
	})
}

func orderUpdateState(db *gorm.DB) {
	ids := []string{"466508742007980032", "466508742041731072"}
	err := db.Select([]string{"state", "updated_at"}).Where("order_id in ?", ids).Updates(&Order{State: 1, UpdatedAt: time.Now()}).Error
	if err != nil {
		logger.Error("queryOrder err", err)
		return
	}
	logger.Info("更新成功")
}

func createManyOrder(db *gorm.DB) {
	m := Order{
		OrderId:        fmt.Sprint(utils.Uint64()),
		UserId:         1,
		Amount:         200,
		Freight:        0,
		DiscountAmount: 0,
		State:          0,
	}
	var src = rand.NewSource(math.MaxInt)
	for i := 0; i < 300; i++ {
		m.OrderId = fmt.Sprint(utils.Uint64())
		m.UserId = int64(rand.New(src).Intn(20)) + 1
		m.Amount = int64(rand.New(src).Intn(15000))
		m.State = rand.New(src).Intn(3)
		err := createOrSaveOrder(db, m)
		if err != nil {
			logger.Error("create order err", err)

		}
	}
}

func createOrSaveOrder(db *gorm.DB, m Order) error {

	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{
			Name: clause.PrimaryKey,
		}},
		OnConstraint: "",
		DoNothing:    false,
		DoUpdates:    clause.Assignments(map[string]interface{}{"state": 1, "updated_at": time.Now()}),
		UpdateAll:    false,
	}).Create(&m).Error
	return err
}
