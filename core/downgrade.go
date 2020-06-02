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
			panic(err)
		}
		db.Commit()
		for _, version := range versions {
			fmt.Printf("DownGrade:%v downgrade success!\n", version)
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
		panic(err.Error)
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
			tables[op.TableName] = 1
			if tableCreateOps[op.TableName] != nil {
				m.createTables([]*Operations{tableCreateOps[op.TableName]})
			} else {
				fmt.Printf("Table: %v not found, downgrade ignored", op.TableName)
			}
		}
	}
	if err := db.DropTableIfExists(t...).Error; err != nil {
		panic(fmt.Sprintf("Delete table failed: %v", err.Error()))
	}
	return tables
}

func (m *Migrate) do(db *gorm.DB, ops *Operations, tableCreateOps map[string]*Operations) {
	if ops != nil && ops.Operations != nil {
		down := m.tableDown(db, ops, tableCreateOps)
		for _, op := range ops.Operations {
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
						panic(fmt.Sprintf("Table: %v DeleteField: %v failed: %v", op.TableName, op.ColumnName, err.Error()))
					}
				}
			case DELETEField:
				if !hasColumn {
					scope := _db.NewScope(op.TableName)
					if err := scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v",
						scope.QuotedTableName(), scope.Quote(op.ColumnName), op.Type)).Exec().DB().Error; err != nil {
						panic(fmt.Sprintf("Table: %v AddField: %v failed: %v", op.TableName, op.ColumnName, err.Error()))
					}
				}
			case ALTERField:
				if hasColumn {
					if err := _db.ModifyColumn(op.ColumnName, op.Type).Error; err != nil {
						panic(fmt.Sprintf("Table: %v AlterField: %v failed: %v", op.TableName, op.ColumnName, err.Error()))
					}
				}
			case ADDIndex:
				if hasIndex {
					if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
						panic(fmt.Sprintf("Table: %v RemoveIndex: %v failed: %v", op.TableName, op.IndexName, err.Error()))
					}
				}
			case ADDUniqueIndex:
				if hasIndex {
					if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
						panic(fmt.Sprintf("Table: %v RemoveUniqueIndex: %v failed: %v", op.TableName, op.IndexName, err.Error()))
					}
				}
			case DELETEIndex:
				if !hasIndex {
					if err := _db.AddIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
						panic(fmt.Sprintf("Table: %v AddIndex: %v failed: %v", op.TableName, op.IndexName, err.Error()))
					}
				}
			case DELETEUniqueIndex:
				if !hasIndex {
					if err := _db.AddUniqueIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
						panic(fmt.Sprintf("Table: %v AddUniqueIndex: %v failed: %v", op.TableName, op.IndexName, err.Error()))
					}
				}
			}
		}
	}
}
