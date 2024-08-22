package database

import (
	"context"
	"fmt"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/logger"
	"gorm.io/gorm"
	"testing"
)

func TestDb_Do(t *testing.T) {
	db, err := NewDb(mysqltool.Config{})
	if err != nil {
		logger.Error(context.Background(), err)
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
	var list []UserPageItem
	db.Do(context.Background(), func(db *gorm.DB) {
		if req.UserId > 0 {
			db = db.Where("user_id=?", req.UserId)
		}
		if req.State > 0 {
			db = db.Where("state=?", req.State)
		}
		if len(req.Name) > 0 {
			db = db.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
		}
		if len(req.Phone) > 0 {
			db = db.Where("phone like ?", fmt.Sprintf("%%%s%%", req.Phone))
		}
		err = db.Count(&total).Order("user_id desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error
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

type UserPageItem struct {
	UserId  int64  `json:"userId"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	State   int    `json:"state"`
	Account string `json:"account"`
}
