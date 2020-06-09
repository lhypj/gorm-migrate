package main

import (
	"dm-gitlab.bolo.me/hubpd/go-migrate/core"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/config"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/migrations"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/models"
)

func main() {
	c := config.GetMIGRATEConfig()
	m := core.Migrate{
		DB:                 models.GetInstance(),
		Migrations:         &migrations.Migrations{},
		Models:             models.GetModels(),
		ModelsRelativePath: "/examples",
	}
	m.Run(c.Command, c.Reversion)
	// MIGRATE_DBDSN="xxx" go run main.go -command makemigrations
	// MIGRATE_DBDSN="xxx" go run main.go -command migrate
	// MIGRATE_DBDSN="xxx" go run main.go -command downgrade -m-reversion 0001_2020060518042059962000
	// MIGRATE_DBDSN="xxx" go run main.go -command fake -m-reversion 0001_2020060518042059962000
	//

}
