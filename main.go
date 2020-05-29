package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"migrate/core"
	"migrate/models"
	"migrate/models/migrations"
)

func main() {

	db, err := gorm.Open("mysql", "root:zxcvbnm123@tcp(localhost:3306)/gomigrate?parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	//db.AutoMigrate()
	//todo use conf
	migrate := core.Migrate{DB: db, ModelsRelativePath: "/models", PackagePath: "migrate/core"}
	migrate.MigrationsInit()
	migrate.MakeMigrations(&migrations.Migrations{}, &models.CreateTableTest{})

	migrate.Migrate(&migrations.Migrations{})
}
