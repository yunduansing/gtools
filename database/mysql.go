package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host     string //mysql host
	Port     int    //mysql port
	Username string //mysql user
	Password string //mysql pwd
	DbName   string
}

// NewMySQL 创建gorm mysql DB
func NewMySQL(c *Config) (*gorm.DB, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

// NewMySQLByDS 创建gorm mysql DB
func NewMySQLByDS(ds string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(ds), &gorm.Config{})
	return db, err
}
