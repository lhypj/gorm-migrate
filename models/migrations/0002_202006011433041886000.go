package migrations

import "migrate/core"

func (*Migrations)Migration_0002_202006011433041886000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.ADDField, TableName: "create_table_tests", ColumnName: "bt", Type: "boolean DEFAULT 0"},
		&core.Operation{Action: core.ALTERField, TableName: "create_table_tests", ColumnName: "name", Type: "varchar(101)", TypeNew: "varchar(50)"},
	)
	return &core.Operations{Revision: "0002_202006011433041886000", DownRevision: []string{"0001_202006011414281206000"}, Operations: ops}
}