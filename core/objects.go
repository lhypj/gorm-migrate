package core

import "github.com/jinzhu/gorm"

type GOrmMigrations struct {
	gorm.Model
	Name string
}

type Migrate struct {
	DB                 *gorm.DB
	PackagePath        string
	ModelsRelativePath string
}

type Field struct {
	Name string
	Type string
}

type Table struct {
	Name   string
	Fields []*Field
}
