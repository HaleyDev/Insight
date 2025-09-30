package data

import (
	c "insight/config"
	"sync"
)

var once sync.Once

func InitData() {
	once.Do(func() {
		if c.GetConfig().MySQL.Enable {
			initMysql()
		}
	})
}
