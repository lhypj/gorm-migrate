package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func (m *Migrate) reversionOps() map[string]*Operations {
	ret := make(map[string]*Operations)
	for _, ops := range m.GetOperations() {
		ret[ops.Revision] = ops
	}
	return ret
}

func reverse(a []string) []string {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func reverseOperation(a []*Operation) []*Operation {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func (m *Migrate) downGradeReversions(version string) []string {
	ao := m.AppliedOrdered()
	for i, reversion := range ao {
		if reversion == version {
			return reverse(ao[i:])
		}
	}
	return []string{}
}

func (m *Migrate) DownGrade(version string) {
	if m.Applied()[version] == 0 {
		fmt.Printf("Version: %v has not applied", version)
		return
	}
	versions := m.downGradeReversions(version)
	if len(versions) == 0 {
		fmt.Println("No versions for downgrade")
		return
	}

	rvOps := m.reversionOps()
	db := m.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			db.Rollback()
			if x, ok := err.(error); ok {
				fmt.Println(x.Error())
			} else {
				fmt.Println(err)
			}
		} else {
			db.Commit()
			for _, version := range versions {
				fmt.Printf("DownGrade:%v downgrade success!\n", version)
			}
		}
	}()

	tableCreateOps := make(map[string]*Operations)
	for reversion, ops := range rvOps {
		if reversion == version {
			break
		}
		for _, op := range ops.Operations {
			if op.Action == ADDTable {
				tableCreateOps[op.TableName] = ops
			}
		}
	}
	for _, vs := range versions {
		m.do(db, rvOps[vs], tableCreateOps)
	}
	if err := db.Unscoped().Where("name in (?)", versions).Delete(&OrmMigrations{}); err != nil {
		panic(err)
	}
}

func (m *Migrate) tableDown(db *gorm.DB, ops *Operations, tableCreateOps map[string]*Operations) map[string]int {
	tables := make(map[string]int)
	var t []interface{}
	for _, op := range ops.Operations {
		if op.Action == ADDTable {
			tables[op.TableName] = 1
			t = append(t, op.TableName)
		}
		if op.Action == DELETETable {
			if tableCreateOps[op.TableName] != nil {
				m.createTables([]*Operations{tableCreateOps[op.TableName]})
			} else {
				fmt.Printf("Table: %v CREATE TABLE INFO NOT FOUND, DOWNGRADE Ignored", op.TableName)
			}
		}
	}
	if err := db.DropTableIfExists(t...).Error; err != nil {
		panic(err)
	}
	return tables
}

func (m *Migrate) do(db *gorm.DB, ops *Operations, tableCreateOps map[string]*Operations) {
	if ops != nil && ops.Operations != nil {
		down := m.tableDown(db, ops, tableCreateOps)
		for _, op := range reverseOperation(ops.Operations) {
			tableName := op.TableName
			if down[tableName] != 0 {
				continue
			}
			_db := db.Table(tableName)
			hasIndex := _db.Dialect().HasIndex(tableName, op.IndexName)
			hasColumn := _db.Dialect().HasColumn(tableName, op.ColumnName)
			switch op.Action {
			case ADDField:
				if hasColumn {
					if err := _db.DropColumn(op.ColumnName).Error; err != nil {
						panic(err)
					}
				}
			case DELETEField:
				if !hasColumn {
					scope := _db.NewScope(op.TableName)
					if err := scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v",
						scope.QuotedTableName(), scope.Quote(op.ColumnName), op.Type)).Exec().DB().Error; err != nil {
						panic(err)
					}
				}
			case ALTERField:
				if hasColumn {
					if err := _db.ModifyColumn(op.ColumnName, op.Type).Error; err != nil {
						panic(err)
					}
				}
			case ADDIndex:
				if hasIndex {
					if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
						panic(err)
					}
				}
			case ADDUniqueIndex:
				if hasIndex {
					if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
						panic(err)
					}
				}
			case DELETEIndex:
				if !hasIndex {
					if err := _db.AddIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
						panic(err)
					}
				}
			case DELETEUniqueIndex:
				if !hasIndex {
					if err := _db.AddUniqueIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
