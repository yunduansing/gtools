package model

import (
	context2 "context"
	"github.com/yunduansing/gtools/database"
	"github.com/yunduansing/gtools/database/pg"
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
		db, err := database.NewDbFromPostgres(pg.Config{Dsn: "host=192.168.2.46 user=postgres password=123456 dbname=user port=5432 sslmode=disable TimeZone=Asia/Shanghai"})
		if err != nil {
			logger.GetLogger().Panic(context2.Background(), "create db context error：", err)
			panic(err)
		}
		dbContext = &DbContext{
			Read: db,
		}
	})
}
