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

func (m *Migrate) Applied() map[string]int {
	ret := make(map[string]int)
	if rows, err := m.DB.Model(&OrmMigrations{}).Select("name").Order("id").Rows(); err != nil {
		panic(err)
	} else {
		var name string
		for rows.Next() {
			if err := rows.Scan(&name); err != nil {
				panic(err)
			} else {
				ret[name] = 1
			}
		}
		return ret
	}
}

func searchUnApplied(node *OperationsNode, Applied map[string]int, unApplied *[]*Operations, allOperations *[]*Operations) {
	if node == nil {
		return
	}
	if !node.IsRoot() {
		if node.Ops == nil {
			panic(fmt.Sprintf("searchUnApplied Failed: some tree node has no Operations"))
		}
		if Applied[node.Ops.Revision] != APPLIED {
			*unApplied = append(*unApplied, node.Ops)
		}
		*allOperations = append(*allOperations, node.Ops)
	}
	for _, child := range node.Children {
		searchUnApplied(child, Applied, unApplied, allOperations)
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
	applied := m.Applied()
	searchUnApplied(root, applied, &unApplied, &allOperations)

	if !checkUnApplied(unApplied, allOperations) {
		panic("UnApplied migrations is not continuous, use Fake or manual handel")
	}
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

func (m *Migrate) createTableAndReturnHandledTable(unApplied []*Operations) map[string]int {
	tableAdded := make(map[string]int)
	tableOps := make(map[string][]*Operation)
	for _, operations := range unApplied {
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

	for tableName := range tableOps {
		scope := m.DB.Table(tableName).NewScope(tableName)
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
		fmt.Println("No UnApplied migrations need to migrate")
		return
	}
	var migrationInfo []string
	repeated := make(map[string]int)

	db := m.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			db.Rollback()
			panic(err)
		}

		db.Commit()
		for _, info := range migrationInfo {
			fmt.Printf("Migrations:%v migrate successful!\n", info)
		}
	}()

	migrated := m.createTableAndReturnHandledTable(unApplied)
	for _, operations := range unApplied {
		repeated[operations.Revision] += 1
		if repeated[operations.Revision] == 1 {
			migrationInfo = append(migrationInfo, operations.Revision)
		}
		for _, op := range operations.Operations {
			if migrated[op.TableName] > 0 {
				continue
			}
			_db := db.Table(op.TableName)
			switch op.Action {
			case ADDField:
				scope := _db.NewScope(op.TableName)
				if err := scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v",
					scope.QuotedTableName(), scope.Quote(op.ColumnName), op.Type)).Exec().DB().Error; err != nil {
					panic(fmt.Sprintf("Table: %v AddField: %v failed: %v", op.TableName, op.ColumnName, err))
				}
			case DELETEField:
				if err := _db.DropColumn(op.ColumnName).Error; err != nil {
					panic(fmt.Sprintf("Table: %v DeleteField: %v failed: %v", op.TableName, op.ColumnName, err))
				}
			case ALTERField:
				if err := _db.ModifyColumn(op.ColumnName, op.TypeNew).Error; err != nil {
					panic(fmt.Sprintf("Table: %v AlterField: %v failed: %v", op.TableName, op.ColumnName, err))
				}
			case ADDIndex:
				if err := _db.AddIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
					panic(fmt.Sprintf("Table: %v AddIndex: %v failed: %v", op.TableName, op.IndexName, err))
				}
			case ADDUniqueIndex:
				if err := _db.AddUniqueIndex(op.IndexName, op.IndexFieldNames...).Error; err != nil {
					panic(fmt.Sprintf("Table: %v AddUniqueIndex: %v failed: %v", op.TableName, op.IndexName, err))
				}
			case DELETETable:
				if err := _db.DropTableIfExists(op.TableName).Error; err != nil {
					panic(fmt.Sprintf("Deletetable: %v failed: %v", op.TableName, err))
				}
			case DELETEIndex:
				if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
					panic(fmt.Sprintf("Table: %v RemoveIndex: %v failed: %v", op.TableName, op.IndexName, err))
				}
			case DELETEUniqueIndex:
				if err := _db.RemoveIndex(op.IndexName).Error; err != nil {
					panic(fmt.Sprintf("Table: %v RemoveUniqueIndex: %v failed: %v", op.TableName, op.IndexName, err))
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

func (m *Migrate) DownMigrate(migrations interface{}) {
	fmt.Println("done")
}
