package core

import (
	"fmt"
	"reflect"
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
			rows.Scan(&name)
			ret[name] = 1
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

func (m *Migrate) UnApplied(migrations interface{}) []*Operations {
	var unApplied []*Operations
	var allOperations []*Operations
	root := m.GetOperationsTree(migrations)
	applied := m.Applied()
	searchUnApplied(root, applied, &unApplied, &allOperations)

	if !checkUnApplied(unApplied, allOperations) {
		panic("UnApplied migrations is not continuous, use Fake or manual handel")
	}
	return unApplied
}

func (*Migrate) Fake(target []string) {

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
		// scope.getTableOptions()?
		s := fmt.Sprintf("CREATE TABLE %v (%v %v)", scope.QuotedTableName(), strings.Join(tags, ","), primaryKeyStr)
		if scope.Raw(s).Exec().HasError() {
			panic(fmt.Sprintf("%v Failed", s))
		}
	}
	return tableAdded
}

func (m *Migrate) Migrate(migrations interface{}) {
	unApplied := m.UnApplied(migrations)
	if len(unApplied) == 0 {
		fmt.Println("No UnApplied migrations need to migrate")
		return
	}
	var migrationInfo []string
	db := m.DB.Begin()
	migrated := m.createTableAndReturnHandledTable(unApplied)
	for _, operations := range m.UnApplied(migrations) {
		migrationInfo = append(migrationInfo, operations.Revision)
		// todo panic err sql to stop 
		for _, op := range operations.Operations {
			if migrated[op.TableName] > 0 {
				continue
			}
			_db := db.Table(op.TableName)
			switch op.Action {
			case ADDField:
				scope := _db.NewScope(op.TableName)
				scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v",
					scope.QuotedTableName(), scope.Quote(op.ColumnName), op.Type)).Exec()
			case DELETEField:
				_db.DropColumn(op.ColumnName)
			case ALTERField:
				_db.ModifyColumn(op.ColumnName, op.TypeNew)
			case ADDIndex:
				fmt.Println("todo: ADDIndex")
			case ADDUniqueIndex:
				fmt.Println("todo: ADDUniqueIndex")
			case DELETETable:
				_db.DropTableIfExists(op.TableName)
			case DELETEIndex:
				fmt.Println("todo: DELETEIndex")
			}
		}
	}

	for _, info := range migrationInfo {
		db.Create(&OrmMigrations{Name: info})
	}
	db.Commit()

	for _, info := range migrationInfo {
		fmt.Printf("Migrations:%v migrate successful!\n", info)
	}
}

func (m *Migrate) DownMigrate(migrations interface{}) {
	fmt.Println("done")
}
