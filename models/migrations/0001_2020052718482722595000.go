package migrations

import "migrate/core"

func (*Migrations)Migration_0001_2020052718060846904000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDTable, TableName: "create_table_tests"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "ID", Type: "int unsigned AUTO_INCREMENT"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "CreatedAt", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "UpdatedAt", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "DeletedAt", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "Name", Type: "varchar(100)"},
		&core.Operation{Action: core.ADDTable, TableName: "create_table_test_v2"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "ID", Type: "int unsigned AUTO_INCREMENT"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "CreatedAt", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "UpdatedAt", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "DeletedAt", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "Name", Type: "varchar(100)"},
	)
	return &core.Operations{Revision: "0001_2020052718482722595000", DownRevision: []string{""}, Operations: ops}
}