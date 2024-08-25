package database

import (
	"context"
	"errors"
	"fmt"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/opentelemetry/tracing"
	"github.com/yunduansing/gtools/utils"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"testing"
	"time"
)

func TestDb_Do(t *testing.T) {
	tracing.InitOtelTracer("db_test")
	ctx := context.WithValue(context.Background(), "requestId", utils.UUID())
	db, err := NewDb(mysqltool.Config{})
	if err != nil {
		logger.Error(ctx, err)
		panic(err)
	}
	req := UserPageReq{
		UserId:   0,
		Name:     "",
		Phone:    "",
		State:    0,
		PageSize: 0,
		Page:     0,
	}
	var total int64
	var list []User
	db.Do(ctx, func(tx *gorm.DB, span trace.Span) *gorm.DB {
		if req.UserId > 0 {
			tx = tx.Where("user_id=?", req.UserId)
		}
		if req.State > 0 {
			tx = tx.Where("state=?", req.State)
		}
		if len(req.Name) > 0 {
			tx = tx.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
		}
		if len(req.Phone) > 0 {
			tx = tx.Where("phone like ?", fmt.Sprintf("%%%s%%", req.Phone))
		}
		result := tx.Count(&total).Order("user_id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list)
		return result
	})
}

type UserPageReq struct {
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	State    int    `json:"state"`
	PageSize int    `json:"pageSize"`
	Page     int    `json:"page"`
}

type User struct {
	UserId   int64  `json:"userId"`
	Username string `json:"name"`
	Phone    string `json:"phone"`
	//State    int    `json:"state"`
	Account string `json:"account"`
}

type UserVip struct {
	Id           int64 `json:"id"`
	UserId       int64 `json:"userId"`
	VipProductId int64 `json:"vipProductId"`
	StartTime    int64 `json:"startTime"` //生效时间
	EndTime      int64 `json:"endTime"`   //失效时间
	State        int   `json:"state"`     //1-购买成功、2-购买失败
}

func TestDb_Find(t *testing.T) {
	tracing.InitOtelTracer("db_test")
	db, err := NewDb(mysqltool.Config{
		Host:     "192.168.1.23",
		Port:     3309,
		Username: "",
		Password: "",
		DbName:   "",
		MaxConn:  0,
		IdleConn: 0,
		LogFile:  "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	var req = struct {
		UserId   int64  `json:"userId"`
		Name     string `json:"name"`
		IsVip    int    `json:"isVip"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	}{
		UserId:   0,
		Name:     "",
		IsVip:    0,
		Page:     1,
		PageSize: 10,
	}
	var users []User
	var count int64
	err = db.Find(context.Background(), &users, func(tx *gorm.DB, span trace.Span) *gorm.DB {
		tx = tx.Table("t_app_user a").Joins("left join t_user_vip b on a.user_id=b.user_id")
		if req.UserId > 0 {
			tx = tx.Where("a.user_id=?", req.UserId)
		}
		if req.Name != "" {
			tx = tx.Where("a.username like ?", fmt.Sprintf("%%%s%%", req.Name))
		}
		if req.IsVip > 0 {
			tx = tx.Where("a.is_vip=?", req.IsVip)
		}
		return tx.Count(&count).Order("user_id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	}).Error
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(count, users)
}

func TestDb_Transaction(t *testing.T) {
	tracing.InitOtelTracer("db_test")
	db, err := NewDb(mysqltool.Config{
		Host:     "192.168.1.23",
		Port:     3309,
		Username: "",
		Password: "",
		DbName:   "",
		MaxConn:  0,
		IdleConn: 0,
		LogFile:  "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	var req = struct {
		UserId       int64  `json:"userId"`
		Name         string `json:"name"`
		VipProductId int64  `json:"vipProductId"`
	}{
		UserId:       100000001,
		Name:         "",
		VipProductId: 100001,
	}

	err = db.Transaction(ctx, func(tx *gorm.DB, span trace.Span) error {
		c := context.WithValue(ctx, "traceId", span.SpanContext().TraceID().String())
		var existsUserVip UserVip
		err = tx.Where("user_id=? and state=1", req.UserId).Order("end_time desc").First(&existsUserVip).Error
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			logger.Error(c, "db find user vip fail", err)
			return err
		}
		start := time.Now()

		newUserVip := UserVip{
			UserId:       req.UserId,
			VipProductId: req.VipProductId,
			StartTime:    start.Unix(),
			//EndTime:      now.Unix(),
			State: 1,
		}
		//y, m, d := now.Date()
		if existsUserVip.Id > 0 && existsUserVip.EndTime >= start.Unix() {
			newUserVip.StartTime = existsUserVip.EndTime
		}
		//假设买的1个月会员
		start = time.Unix(newUserVip.StartTime, 0).AddDate(0, 1, 0)
		newUserVip.EndTime = start.Unix()
		err = tx.Create(&newUserVip).Error
		if err != nil {
			logger.Error(c, "db create user vip fail", err)
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDb_Save(t *testing.T) {
	tracing.InitOtelTracer("db_test")
	db, err := NewDb(mysqltool.Config{
		Host:     "192.168.1.23",
		Port:     3309,
		Username: "",
		Password: "",
		DbName:   "",
		MaxConn:  0,
		IdleConn: 0,
		LogFile:  "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	var req = struct {
		UserId int64  `json:"userId"`
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		State  int    `json:"state"`
	}{
		UserId: 0,
		Name:   "Bob",
		Phone:  "13311112222",
		State:  1,
	}

	var newUser = User{
		UserId:   req.UserId,
		Username: "",
		Phone:    req.Phone,
		Account:  req.Phone,
	}

	err = db.Save(ctx, &newUser, clause.OnConflict{
		Columns:   []clause.Column{{Name: "phone"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"state": 1}),
	}).Error
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(newUser)
}
