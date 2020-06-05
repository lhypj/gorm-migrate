package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"dm-gitlab.bolo.me/hubpd/go-migrate/core"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/migrate/migrations"
	"dm-gitlab.bolo.me/hubpd/go-migrate/examples/models"
)

func main() {
	// 可以用自己搞的db
	db, err := gorm.Open("mysql", "root:zxcvbnm123@tcp(localhost:3306)/gomigrate?parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// todo 自动捞model ???
	mds := []interface{}{&models.CreateTableTest{}, &models.CreateTableTestV2{}}

	m := core.Migrate{
		DB:                 db,
		Migrations:         &migrations.Migrations{},
		Models:             mds,
	}
	m.Run()
}
