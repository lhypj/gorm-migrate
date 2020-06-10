package main

import (
	migrate "github.com/lhypj/gorm-migrate/core"
	"github.com/lhypj/gorm-migrate/examples/config"
	"github.com/lhypj/gorm-migrate/examples/migrations"
	"github.com/lhypj/gorm-migrate/examples/models"
)


func main() {

	c := config.GetConfig()
	migrate := migrate.Migrate{
		DB: models.GetInstance(),
		//ModelsRelativePath: "/examples",
		Migrations: &migrations.Migrations{},
		Models: []interface{}{
			&models.CreateTableTest{},
			&models.CreateTableTestV2{},
		},
	}
	migrate.MigrationsInit()
	migrate.Run(c.Command, c.Reversion)
	//
	//migrate.Migrate()
	//migrate.Fake("0005_202006011508224506000")
	//migrate.Merge()
	//migrate.List()
	//migrate.DownGrade("0005_202006011508224506000")
}

