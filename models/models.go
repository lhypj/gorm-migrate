package models

import (
	"github.com/jinzhu/gorm"
)

type CreateTableTest struct {
	gorm.Model
	Name string `gorm:"type:varchar(101)"`
}

type CreateTableTestV2 struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}