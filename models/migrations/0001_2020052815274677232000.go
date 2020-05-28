package migrations

import "migrate/core"

func (*Migrations)Migration_0001_2020052718060846904000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDTable, TableName: "create_table_tests"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "deleted_at", Type: "DATETIME NULL", TypeNew: "fff"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "name", Type: "varchar(100)"},
	)
	return &core.Operations{Revision: "0001_2020052815274677232000", DownRevision: []string{""}, Operations: ops}
}