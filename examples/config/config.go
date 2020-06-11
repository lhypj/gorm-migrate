package config

import (
	"github.com/koding/multiconfig"
	"sync"
)


type Command struct {
	GrpcServerAddress string
	ConsumerGroupID   string
	SentryDsn         string
	KafkaHosts        []string
	DBDsn             string

	Migrate   bool
	Command   string
	Reversion string
}

var once sync.Once
var c *Command

func GetConfig() *Command {
	once.Do(func() {
		c = new(Command)
		m := multiconfig.New()
		m.MustLoad(c)
	})
	return c
}

func init() {
	GetConfig()
}
