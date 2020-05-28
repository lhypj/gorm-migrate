package core

import "github.com/jinzhu/gorm"

type OrmMigrations struct {
	gorm.Model
	Name string `gorm:"unique;not null;size:30"`
}

type Migrate struct {
	DB                 *gorm.DB
	PackagePath        string
	ModelsRelativePath string
}

type Field struct {
	Name             string
	Type             string
	IsPrimary        bool
	IndexNames       []string
	UniqueIndexNames []string
}

type Table struct {
	Name   string
	Fields []*Field
}
