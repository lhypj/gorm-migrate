package config

import (
	"dm-gitlab.bolo.me/hubpd/go-migrate/core"
	"github.com/koding/multiconfig"
	"sync"
)

type Hubpd struct {
	core.M
	DBDsn string
}

var once sync.Once
var c *Hubpd

func GetConfig() *Hubpd {
	once.Do(func() {
		c = new(Hubpd)
		m := multiconfig.New()
		m.MustLoad(c)
	})
	return c
}

func init() {
	GetConfig()
}
