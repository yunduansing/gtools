package config

import (
	"context"
	"github.com/yunduansing/gtools/database"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/examples/coupon/model"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/yaml"
)

var (
	writeMysqlConfig mysqltool.Config
	readMysqlConfig  mysqltool.Config
	logConfig        logger.Config
	app              App
)

var (
	Redis *redis.Client
)

type App struct {
	Name  string
	Port  int
	Mysql mysqltool.Config
	Redis redis.Config
	Log   logger.Config
}

func loadConfig() error {
	if err := yaml.Resolver("etc/app.yaml", &app); err != nil {
		return err
	}
	return nil
}

func InitConfig() {
	if err := loadConfig(); err != nil {
		panic(err)
	}

	logger.InitLog(logConfig)

	if err := database.InitDb(model.WriteMysql, writeMysqlConfig); err != nil {
		logger.GetLogger().Error(context.TODO(), "Init WriteMysql err:", err)
		panic(err)
	}
	if err := database.InitDb(model.ReadMysql, readMysqlConfig); err != nil {
		logger.GetLogger().Error(context.TODO(), "Init ReadMysql err:", err)
		panic(err)
	}

	initRedis()
}

func initRedis() {
	Redis = redis.New(app.Redis)
}
