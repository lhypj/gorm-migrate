package migrations

import "migrate/core"

func (*Migrations) Migration_0006_2020060117392545970000() *core.Operations {
	return &core.Operations{
		Operations: []*core.Operation{},
		Revision:   "0006_2020060117392545970000",
		DownRevision: []string{
			"0005_202006011508224506000", 
			"0005_202006011508224506333",
		},
	}
}