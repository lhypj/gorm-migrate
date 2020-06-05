package migrations

import "dm-gitlab.bolo.me/hubpd/go-migrate/core"

func (*Migrations) Migration_0001_2020060518042059962000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDTable, TableName: "create_table_test_v2"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test_v2", ColumnName: "name", Type: "varchar(100)"},
		&core.Operation{TableName: "create_table_test_v2", Action: core.ADDIndex, IndexName: "idx_create_table_test_v2_deleted_at", IndexFieldNames: []string{"deleted_at"}},
		&core.Operation{Action: core.ADDTable, TableName: "create_table_test"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "name", Type: "varchar(50)"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "bt", Type: "boolean DEFAULT 0"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_test", ColumnName: "uidx1", Type: "varchar(50)"},
		&core.Operation{TableName: "create_table_test", Action: core.ADDIndex, IndexName: "idx_create_table_test_deleted_at", IndexFieldNames: []string{"deleted_at"}},
		&core.Operation{TableName: "create_table_test", Action: core.ADDIndex, IndexName: "idx_create_table_test_bt", IndexFieldNames: []string{"bt"}},
		&core.Operation{TableName: "create_table_test", Action: core.ADDUniqueIndex, IndexName: "name_u_idx", IndexFieldNames: []string{"name"}},
		&core.Operation{TableName: "create_table_test", Action: core.ADDUniqueIndex, IndexName: "uidx", IndexFieldNames: []string{"uidx1"}},
		&core.Operation{Action: core.ADDTable, TableName: "for"},
		&core.Operation{Action: core.ADDField, TableName: "for", ColumnName: "id", Type: "int unsigned AUTO_INCREMENT", IsPrimary: true},
		&core.Operation{Action: core.ADDField, TableName: "for", ColumnName: "created_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "for", ColumnName: "updated_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "for", ColumnName: "deleted_at", Type: "DATETIME NULL"},
		&core.Operation{Action: core.ADDField, TableName: "for", ColumnName: "nn", Type: "varchar(100)"},
		&core.Operation{TableName: "for", Action: core.ADDIndex, IndexName: "idx_for_deleted_at", IndexFieldNames: []string{"deleted_at"}},
	)
	return &core.Operations{Revision: "0001_2020060518042059962000", DownRevision: []string{""}, Operations: ops}
}