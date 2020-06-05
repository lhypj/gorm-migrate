package migrations

import "dm-gitlab.bolo.me/hubpd/gorm-migrate/core"

func (*Migrations) Migration_0001_202006041853553761000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDTable, TableName: "create_table_tests"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "name", Type: "varchar(50)"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "bt", Type: "boolean DEFAULT 0"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "uidx1", Type: "varchar(50)"},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDIndex, IndexName: "idx_create_table_tests_deleted_at", IndexFieldNames: []string{"deleted_at"}},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDIndex, IndexName: "idx_create_table_tests_bt", IndexFieldNames: []string{"bt"}},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDUniqueIndex, IndexName: "name_u_idx", IndexFieldNames: []string{"name"}},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDUniqueIndex, IndexName: "uidx", IndexFieldNames: []string{"uidx1"}},
		&core.Operation{Action: core.ADDTable, TableName: "create_table_test_v2"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "name", Type: "varchar(100)"},
		&core.Operation{TableName: "create_table_test_v2", Action: core.ADDIndex, IndexName: "idx_create_table_test_v2_deleted_at", IndexFieldNames: []string{"deleted_at"}},
	)
	return &core.Operations{Revision: "0001_202006041853553761000", DownRevision: []string{""}, Operations: ops}
}