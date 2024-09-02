package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres create postgres conn
//
// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
func NewPostgres(c *Config) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(c.Dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
