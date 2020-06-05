package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"sort"
	"strings"
)

const (
	APPLIED = 1
)

func (m *Migrate) AppliedOrdered() []string {
	var ret []string
	if rows, err := m.DB.Model(&OrmMigrations{}).Select("name").Order("id").Rows(); err != nil {
		panic(err)
	} else {
		var name string
		for rows.Next() {
			if err := rows.Scan(&name); err != nil {
				panic(err)
			} else {
				ret = append(ret, name)
			}
		}
	}
	return ret
}

func (m *Migrate) Applied() map[string]int {
	ret := make(map[string]int)
	for _, name := range m.AppliedOrdered() {
		ret[name] = 1
	}
	return ret
}

func searchUnApplied(node *OperationsNode, Applied map[string]int, unApplied *[]*Operations, allOperations *[]*Operations) {
	if node == nil {
		return
	}
	var nodes []*OperationsNode
	nodes = append(nodes, node)
	searched := make(map[string]int)

	for len(nodes) != 0 {
		node = nodes[0]
		nodes = nodes[1:]
		if !node.IsRoot() && node.Ops != nil && searched[node.Ops.Revision] == 0 {
			if Applied[node.Ops.Revision] != APPLIED {
				*unApplied = append(*unApplied, node.Ops)
			}
			*allOperations = append(*allOperations, node.Ops)
			searched[node.Ops.Revision] = 1
		}
		for _, child := range node.Children {
			nodes = append(nodes, child)
		}
	}
}

func checkUnApplied(unApplied []*Operations, allOperations []*Operations) bool {
	if len(unApplied) == 0 {
		return true
	}
	return reflect.DeepEqual(unApplied, allOperations[len(allOperations)-len(unApplied):])
}

func (m *Migrate) UnApplied() []*Operations {
	var unApplied []*Operations
	var allOperations []*Operations
	root := m.GetOperationsTree(true)
	searchUnApplied(root, m.Applied(), &unApplied, &allOperations)
	sort.Sort(OperationSlice(unApplied))
	return unApplied
}

func (*Migrate) getTableOptions(scope *gorm.Scope) string {
	tableOptions, ok := scope.Get("gorm:table_options")
	if !ok {
		return ""
	}
	return " " + tableOptions.(string)
}

func (m *Migrate) tableOpsForCreate(ops []*Operations) map[string][]*Operation {
	tableOps := make(map[string][]*Operation)
	for _, operations := range ops {
		for _, op := range operations.Operations {
			if op.Action == ADDTable {
				tableOps[op.TableName] = append(tableOps[op.TableName], op)
			} else {
				if len(tableOps[op.TableName]) != 0 {
					tableOps[op.TableName] = append(tableOps[op.TableName], op)
				}
			}
		}
	}
	return tableOps
}

func (m *Migrate) createTables(unApplied []*Operations) map[string]int {
	tableAdded := make(map[string]int)
	tableOps := m.tableOpsForCreate(unApplied)

	for tableName := range tableOps {
		scope := m.DB.Table(tableName).NewScope(tableName)
		if scope.Dialect().HasTable(tableName) {
			continue
		}
		var tags []string
		var primaryKeys []string
		var primaryKeyInColumnType = false
		tableAdded[tableName] = 1
		for _, op := range tableOps[tableName] {
			// todo 如支持外键，则需考虑此op是本model还是外部model 暂不支持外键
			if strings.Contains(strings.ToLower(op.Type), "primary key") {
				primaryKeyInColumnType = true
			}
			if op.IsPrimary {
				primaryKeys = append(primaryKeys, scope.Quote(op.ColumnName))
			}
			if op.ColumnName != "" && op.Type != "" {
				tags = append(tags, scope.Quote(op.ColumnName)+" "+op.Type)
			}
		}
		var primaryKeyStr string
		if len(primaryKeys) > 0 && !primaryKeyInColumnType {
			primaryKeyStr = fmt.Sprintf(", PRIMARY KEY (%v)", strings.Join(primaryKeys, ","))
		}
		s := fmt.Sprintf("CREATE TABLE %v (%v %v)%s", scope.QuotedTableName(), strings.Join(tags, ","), primaryKeyStr, m.getTableOptions(scope))
		if scope.Raw(s).Exec().HasError() {
			panic(fmt.Sprintf("%v Failed", s))
		}
	}
	return tableAdded
}

func (m *Migrate) Migrate() {
	unApplied := m.UnApplied()
	if len(unApplied) == 0 {
		fmt.Println("No unApplied migrations need to migrate")
		return
	}
	var migrationInfo []string

	db := m.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			db.Rollback()
			fmt.Println("Migrate failed!")
			panic(err)
		} else {
			db.Commit()
			for _, info := range migrationInfo {
				fmt.Printf("Migrations:%v migrate success!\n", info)
			}
		}
	}()

	migrated := m.createTables(unApplied)
	for _, operations := range unApplied {
		migrationInfo = append(migrationInfo, operations.Revision)
		for _, op := range operations.Operations {
			tableCreated := migrated[op.TableName] > 0
			_db := db.Table(op.TableName)
			hasIndex := _db.Dialect().HasIndex(op.TableName, op.IndexName)
			hasColumn := _db.Dialect().HasColumn(op.TableName, op.ColumnName)
			switch op.Action {
			case ADDField:
				if !hasColumn && !tableCreated {
					scope := _db.NewScope(op.TableName)
					if err := scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v",
						scope.QuotedTableName(), scope.Quote(op.ColumnName), op.Type)).Exec().DB().Error; err != nil {
						panic(err)
					}
				}
			case DELETEField:
				if hasColumn {
					if err := _db.DropColumn(op.ColumnName).Error; err != nil {
						panic(err)
					}
				}
			case ALTERField:
				if hasColumn {
					if err := _db.ModifyColumn(op.ColumnName, op.TypeNew).Error; err != nil {
						panic(err)
					}
				}

			case ADDIndex:
				if !hasIndex {
					if err := _db.AddIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
						panic(err)
					}
				}

			case ADDUniqueIndex:
				if !hasIndex {
					if err := _db.AddUniqueIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
						panic(err)
					}
				}
			case DELETETable:
				if err := _db.DropTableIfExists(op.TableName).Error; err != nil {
					panic(err)
				}
			case DELETEIndex:
				if hasIndex {
					if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
						panic(err)
					}
				}
			case DELETEUniqueIndex:
				if hasIndex {
					if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
						panic(err)
					}
				}
			}
		}
	}

	for _, info := range migrationInfo {
		if err := db.Create(&OrmMigrations{Name: info}).Error; err != nil {
			panic(err)
		}
	}
}
