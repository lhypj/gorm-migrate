package models

import (
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	otgorm "github.com/smacker/opentracing-gorm"
	"sync"
)

var instance *gorm.DB
var once sync.Once

func GetInstance() *gorm.DB {
	once.Do(func() {
		var err error
		dsn := config.GetConfig().DBDsn
		if len(dsn) == 0 {
			log.Fatal("Failed to get database dsn")
		}
		instance, err = gorm.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("open db: %s", err)
		}
		instance.SingularTable(true)
		otgorm.AddGormCallbacks(instance)
	})
	return instance
}

func GetModels () []interface{} {
	return []interface{}{
		&CreateTableTestV2{},
		&CreateTableTest{},
		&For{},
	}
}
