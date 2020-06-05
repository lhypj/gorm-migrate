package models

import (
	"github.com/jinzhu/gorm"
)

type CreateTableTest struct {
	gorm.Model
	Name  string `gorm:"type:varchar(50);unique_index:name_u_idx;"`
	BT    bool   `gorm:"default:0;index"`
	Uidx1 string `gorm:"type:varchar(50);unique_index:uidx;"`
}

type For struct {
	gorm.Model
	NN string `gorm:"type:varchar(100)"`
}

type CreateTableTestV2 struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
	//For For 不支持
}
