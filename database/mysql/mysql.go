package mysqltool

import (
	"fmt"
	sysLog "github.com/yunduansing/gtools/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
)

type Config struct {
	Host     string //mysql host
	Port     int    //mysql port
	Username string //mysql user
	Password string //mysql pwd
	DbName   string //db name
	MaxConn  int    //最大连接数 默认200
	IdleConn int    //空闲时连接数 默认20
	LogFile  string `json:",default=log/db"`
}

type Option func()

func WithMaxIdleConn() {

}

func WithMaxOpenConn() {

}

var (
	idleConn = 20
	maxConn  = 200
)

// NewMySQLFromConfig 创建gorm mysql DB
func NewMySQLFromConfig(c *Config, opts ...Option) (*gorm.DB, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.DbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			LogLevel: logger.Info,
		}),
	})
	for _, opt := range opts {
		opt()
	}
	if err != nil {
		sysLog.Error("open mysql err:", err)
		panic(err)
	}
	ic := idleConn
	mc := maxConn
	if c.IdleConn > 0 {
		ic = c.IdleConn
	}
	if c.MaxConn > 0 {
		mc = c.MaxConn
	}
	db.Use(dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(ic).SetMaxOpenConns(mc).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(24 * time.Hour))

	return db, err
}

func SetMultiDb(db *gorm.DB) {

}

// NewMySQLFromConnString 创建gorm mysql DB
func NewMySQLFromConnString(ds string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(ds), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
			LogLevel: logger.Info,
		}),
	})
	if err != nil {
		sysLog.Error("open mysql err:", err)
		panic(err)
	}
	ic := idleConn
	mc := maxConn
	db.Use(dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(ic).SetMaxOpenConns(mc).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(24 * time.Hour))
	return db, nil
}

// NewMySQL 创建gorm mysql DB
func NewMySQL(c *Config) (*gorm.DB, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		sysLog.Error("open mysql err:", err)
		panic(err)
	}
	ic := idleConn
	mc := maxConn
	db.Use(dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(ic).SetMaxOpenConns(mc).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(24 * time.Hour))
	return db, err
}

// NewMySQLByDS 创建gorm mysql DB
func NewMySQLByDS(ds string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(ds), &gorm.Config{})
	if err != nil {
		sysLog.Error("open mysql err:", err)
		panic(err)
	}
	ic := idleConn
	mc := maxConn
	db.Use(dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(ic).SetMaxOpenConns(mc).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(24 * time.Hour))
	return db, err
}
