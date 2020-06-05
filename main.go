package main

import (
	"dm-gitlab.bolo.me/hubpd/go-migrate/core"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/config"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/migrations"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/models"
)

func main() {
	c := config.GetConfig()
	m := core.Migrate{
		DB:                 models.GetInstance(),
		Migrations:         &migrations.Migrations{},
		Models:             models.GetModels(),
		ModelsRelativePath: "/examples",
	}
	m.Run(c.Command, c.Reversion)

	// HUBPD_DBDSN="xxx" go run main.go -m-command makemigrations
	// HUBPD_DBDSN="xxx" go run main.go -m-command migrate
	// HUBPD_DBDSN="xxx" go run main.go -m-command downgrade -m-reversion 0001_2020060518042059962000
	// HUBPD_DBDSN="xxx" go run main.go -m-command fake -m-reversion 0001_2020060518042059962000
	//
}
