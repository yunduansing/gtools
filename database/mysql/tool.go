package mysqltool

import (
	"gorm.io/gorm"
	"sync"
)

var (
	rwLock   sync.RWMutex
	mysqlMap = make(map[string]*gorm.DB)
)

func InitMysql(key string, c *Config) (err error) {
	rwLock.Lock()
	defer rwLock.Unlock()
	if _, ok := mysqlMap[key]; ok {
		return nil
	}
	mysqlMap[key], err = NewMySQLFromConfig(c)
	return nil
}

func Get(key string) *gorm.DB {
	rwLock.RLock()
	defer rwLock.RUnlock()
	return mysqlMap[key]
}
