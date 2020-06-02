package migrations

import "migrate/core"

func (*Migrations) Migration_0007_2020060211044311701000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{Action: core.DELETEField, TableName: "create_table_tests", ColumnName: "uidx2"},
		&core.Operation{TableName: "create_table_tests", Action: core.DELETEUniqueIndex, IndexName: "uidx", IndexFieldNames: []string{"uidx1", "uidx2"}},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDUniqueIndex, IndexName: "uidx", IndexFieldNames: []string{"uidx1"}},
	)
	return &core.Operations{Revision: "0007_2020060211044311701000", DownRevision: []string{"0006_2020060117392545970000"}, Operations: ops}
}