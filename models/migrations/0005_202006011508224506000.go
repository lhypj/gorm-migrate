package migrations

import "migrate/core"

func (*Migrations) Migration_0005_202006011508224506000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDTable, TableName: "create_table_test_v2"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "name", Type: "varchar(100)"},
		&core.Operation{Action: core.ADDIndex, TableName: "create_table_test_v2", IndexName: "idx_create_table_test_v2_deleted_at", IndexFieldNames: []string{"deleted_at"}},
	)
	return &core.Operations{Revision: "0005_202006011508224506000", DownRevision: []string{"0004_202006011507162482000"}, Operations: ops}
}
