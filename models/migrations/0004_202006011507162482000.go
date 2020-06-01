package migrations

import "migrate/core"

func (*Migrations)Migration_0004_202006011507162482000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "uidx1", Type: "varchar(50)"},
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "uidx2", Type: "varchar(50)"},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDUniqueIndex, IndexName: "uidx", IndexFieldNames: []string{"uidx1", "uidx2"}},
	)
	return &core.Operations{Revision: "0004_202006011507162482000", DownRevision: []string{"0003_2020060114364693063000"}, Operations: ops}
}