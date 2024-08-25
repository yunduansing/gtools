package config

import (
	"context"
	mysqltool "github.com/yunduansing/gtools/database/mysql"
	"github.com/yunduansing/gtools/examples/coupon/model"
	"github.com/yunduansing/gtools/logger"
	"github.com/yunduansing/gtools/redistool"
	"github.com/yunduansing/gtools/yaml"
)

var (
	writeMysqlConfig mysqltool.Config
	readMysqlConfig  mysqltool.Config
	logConfig        logger.Config
	app              App
)

var (
	Redis *redistool.Client
)

type App struct {
	Name  string
	Port  int
	Mysql mysqltool.Config
	Redis redistool.Config
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

	if err := mysqltool.InitMysql(model.WriteMysql, &writeMysqlConfig); err != nil {
		logger.Error(context.TODO(), "Init WriteMysql err:", err)
		panic(err)
	}
	if err := mysqltool.InitMysql(model.ReadMysql, &readMysqlConfig); err != nil {
		logger.Error(context.TODO(), "Init ReadMysql err:", err)
		panic(err)
	}

	initRedis()
}

func initRedis() {
	Redis = redistool.New(app.Redis)
}
