package config

import (
	"github.com/koding/multiconfig"
	"sync"
)

type MigrateLoader struct {
	DBDsn     string `default:"root:zxcvbnm123@tcp(localhost:3306)/TestDb?parseTime=True&loc=Asia%2FShanghai"`
	Command   string
	Reversion string
}

var once sync.Once
var c *MigrateLoader

func GetConfig() *MigrateLoader {
	once.Do(func() {
		c = new(MigrateLoader)
		m := multiconfig.New()
		m.MustLoad(c)
	})
	return c
}

func init() {
	GetConfig()
}
