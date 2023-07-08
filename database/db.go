package database

import "gorm.io/gorm"

type Db struct {
	Mysql *gorm.DB
}

var Context Db
