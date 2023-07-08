package mysqltool

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Config struct {
	Host     string //mysql host
	Port     int    //mysql port
	Username string //mysql user
	Password string //mysql pwd
	DbName   string //db name
	MaxConn  int    //最大连接数
	IdleConn int    //空闲时连接数
	LogFile  string `json:",default=log/db"`
}

var (
	idleConn = 20
	maxConn  = 200
)

// NewMySQLFromConfig 创建gorm mysql DB
func NewMySQLFromConfig(c *Config) (*gorm.DB, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	ll, err := rotatelogs.New(c.LogFile+"-%Y%m%d.log", rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(ll, "", log.LstdFlags), logger.Config{
			LogLevel: logger.Info,
		}),
	})
	if err == nil {
		ic := idleConn
		mc := maxConn
		if c.IdleConn > 0 {
			ic = c.IdleConn
		}
		if c.MaxConn > 0 {
			mc = c.MaxConn
		}
		db, err := db.DB()
		if err == nil {
			db.SetMaxIdleConns(ic)
			db.SetMaxOpenConns(mc)
		} else {
			return nil, err
		}
	}
	return db, err
}

// NewMySQLFromConnString 创建gorm mysql DB
func NewMySQLFromConnString(ds string) (*gorm.DB, error) {
	ll, err := rotatelogs.New("log/db-%Y%m%d.log", rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(ds), &gorm.Config{
		Logger: logger.New(log.New(ll, "", log.LstdFlags), logger.Config{
			LogLevel: logger.Info,
		}),
	})
	if err == nil {
		db, err := db.DB()
		if err == nil {
			db.SetMaxIdleConns(idleConn)
			db.SetMaxOpenConns(maxConn)
		} else {
			return nil, err
		}
	}
	return db, err
}
