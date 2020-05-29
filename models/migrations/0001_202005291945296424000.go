package migrations

import "migrate/core"

func (*Migrations)Migration_0001_202005291945296424000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDTable, TableName: "create_table_tests"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "name", Type: "varchar(101)"},
		&core.Operation{Action: core.ADDIndex, TableName: "create_table_tests", IndexName: "idx_create_table_tests_deleted_at", IndexFieldNames: []string{"deleted_at"}},
	)
	return &core.Operations{Revision: "0001_202005291945296424000", DownRevision: []string{""}, Operations: ops}
}