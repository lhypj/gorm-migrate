package migrations

import "migrate/core"

func (*Migrations)Migration_0003_2020060114364693063000() *core.Operations {
	var ops []*core.Operation
	ops = append(ops,
		&core.Operation{TableName: "create_table_tests", Action: core.ADDIndex, IndexName: "idx_create_table_tests_bt", IndexFieldNames: []string{"bt"}},
		&core.Operation{TableName: "create_table_tests", Action: core.ADDUniqueIndex, IndexName: "name_u_idx", IndexFieldNames: []string{"name"}},
	)
	return &core.Operations{Revision: "0003_2020060114364693063000", DownRevision: []string{"0002_202006011433041886000"}, Operations: ops}
}