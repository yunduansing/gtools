package redistool

import (
	context2 "context"
	"fmt"
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/database"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/utils"
	"sync"
	"testing"
	"time"
)

type Order struct {
	Id          int64
	OrderNumber string
	Amount      int64
	UserId      int64
}

func (*Order) TableName() string {
	return "t_order"
}

func TestLocker_Acquire(t *testing.T) {
	cli := New(Config{Addr: []string{""}, Password: ""})

	key := "gtools:redis:lock:test"

	db, err := database.NewDb(mysqltool.Config{
		Host:     "",
		Port:     3306,
		Username: "",
		Password: "",
		DbName:   "",
	})

	if err != nil {
		t.Error("init db err ", err)
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(index int) {
			ctx := context.NewContext(context2.Background(), context.WithRequestId(fmt.Sprint(index)))
			locker := NewRedisLockWithContext(ctx, cli.UniversalClient, key)
			locker.SetExpire(15)
			getSuccess, err := locker.Acquire(3000, time.Millisecond*10)
			if err != nil {
				t.Error(i, "locker.Acquire err ", index, "   ===>", err)
				return
			}
			if !getSuccess {
				t.Error(i, "locker.Acquire fail ", index, "   ===>")
				return
			}

			defer func() {
				ctx.Log.Infof(ctx.Ctx, "locker.Release success on %d", index)
				locker.Release()

				wg.Done()
			}()

			order := Order{
				Id:          0,
				OrderNumber: utils.NewIDFromUint64(),
				Amount:      1,
				UserId:      1,
			}

			ctx.Log.WithField("order", order).Infof(ctx.Ctx, "locker.Acquire success on %d ", index)

			var userOrderExists int64
			if err = db.DB.Model(&Order{}).Where("user_id = ?", order.UserId).Count(&userOrderExists).Error; err != nil {
				ctx.Log.Error(ctx.Ctx, "count user order err ", err)
				return
			}
			if userOrderExists > 0 {
				return
			}

			<-time.After(time.Millisecond * 10)
			err = db.DB.Create(&order).Error
			if err != nil {
				ctx.Log.Error(ctx.Ctx, "create order err ", err)
				return
			}
		}(i)
	}
	wg.Wait()

}
