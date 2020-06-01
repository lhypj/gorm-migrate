package migrations

import "migrate/core"

func (*Migrations) Migration_0005_202006011508224506333() *core.Operations {
	return &core.Operations{
		Operations: []*core.Operation{},
		Revision: "0005_202006011508224506333",
		DownRevision: []string{
			"0004_202006011507162482000",
		},
	}
}

