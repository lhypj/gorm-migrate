package main

import (
	migrate "github.com/lhypj/gorm-migrate/core"
	"github.com/lhypj/gorm-migrate/examples/config"
	"github.com/lhypj/gorm-migrate/examples/migrations"
	"github.com/lhypj/gorm-migrate/examples/models"
)

func main() {
	c := config.GetConfig()
	db := models.GetInstance()
	if c.Migrate {
		m := migrate.Migrate{
			DB:         db,
			Migrations: &migrations.Migrations{},
			Models: []interface{}{
				&models.CreateTableTest{},
				&models.CreateTableTestV2{},
			},
		}
		m.MigrationsInit()
		m.Run(c.Command, c.Reversion)
		return
	}
	//  COMMAND_DBDSN="xxx" go run main.go -command migrate -migrate
	// other env handle ...
}
