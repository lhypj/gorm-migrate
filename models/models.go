package models

import (
	"github.com/jinzhu/gorm"
)

type CreateTableTest struct {
	gorm.Model
	Name  string `gorm:"type:varchar(50);unique_index:name_u_idx;"`
	BT    bool   `gorm:"default:0;index"`
	Uidx1 string `gorm:"type:varchar(50);unique_index:uidx;"`
	Uidx2 string `gorm:"type:varchar(50);unique_index:uidx;"`
}

type CreateTableTestV2 struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}
