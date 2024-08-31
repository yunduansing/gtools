package model

import (
	context2 "context"
	"github.com/yunduansing/gtools/database"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/logger"
	"sync"
)

type DbContext struct {
	Read *database.Db
}

var (
	dbContext *DbContext
	once      sync.Once
)

func InitDbContext() {
	once.Do(func() {
		db, err := database.NewDb(mysqltool.Config{
			Host:     "192.168.6.23",
			Port:     3309,
			Username: "",
			Password: "",
			DbName:   "",
			MaxConn:  0,
			IdleConn: 0,
			LogFile:  "./logs",
		})
		if err != nil {
			logger.GetLogger().Panic(context2.Background(), "create db context error：", err)
			panic(err)
		}
		dbContext = &DbContext{
			Read: db,
		}
	})
}
