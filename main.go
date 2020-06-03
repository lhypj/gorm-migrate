package main

import (
	"fmt"
	"os"
)

func main() {

	//db, err := gorm.Open("mysql", "root:zxcvbnm123@tcp(localhost:3306)/gomigrate?parseTime=True&loc=Asia%2FShanghai")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer db.Close()
	//db.AutoMigrate()
	//todo use conf
	//migrate := migrate.Migrate{DB: db, ModelsRelativePath: "/models", PackagePath: "migrate/core", Migrations: &migrations.Migrations{}}
	//migrate.MigrationsInit()
	//migrate.MakeMigrations(&models.CreateTableTest{}, &models.CreateTableTestV2{})
	//
	//migrate.Migrate()
	//migrate.Fake("0005_202006011508224506000")
	//migrate.Merge()
	//migrate.List()
	//migrate.DownGrade("0005_202006011508224506000")

	var migrationVersion string
	if len(os.Args) < 2 {
		fmt.Println("need command")
		return
	}
	command := os.Args[1]
	if len(os.Args) > 2 {
		migrationVersion = os.Args[2]
	}
	fmt.Println(command, migrationVersion)
}
