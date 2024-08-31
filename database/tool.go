package database

import (
	"github.com/yunduansing/gtools/database/mysql"
	"sync"
)

var (
	rwLock sync.RWMutex
	dbMap  = make(map[string]*Db)
)

func InitDb(key string, c mysqltool.Config) (err error) {
	rwLock.Lock()
	defer rwLock.Unlock()
	if _, ok := dbMap[key]; ok {
		return nil
	}
	dbMap[key], err = NewDb(c)
	return nil
}

func Get(key string) *Db {
	rwLock.RLock()
	defer rwLock.RUnlock()
	return dbMap[key]
}
