package config

import (
	"dm-gitlab.bolo.me/hubpd/go-migrate/core"
	"github.com/koding/multiconfig"
	"sync"
)

func GetMIGRATEConfig() *core.MIGRATE {
	var once sync.Once
	var c *core.MIGRATE
	once.Do(func() {
		c = new(core.MIGRATE)
		m := multiconfig.New()
		m.MustLoad(c)
	})
	return c
}

func init() {
	GetMIGRATEConfig()
}
