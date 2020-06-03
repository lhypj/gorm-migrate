package core

import (
	"fmt"
	"strconv"
	"strings"
)

func maxHead(heads []string) string {
	var iMax int
	r := make(map[int]string)
	for _, head := range heads {
		h := strings.Split(head, MIGRATIONSplit)[0]
		if i, err := strconv.Atoi(h); err != nil {
			panic(err)
		} else {
			if i > iMax {
				iMax = i
			}
			r[i] = head
		}
	}
	if val, ok := r[iMax]; ok {
		return val
	} else {
		panic("Internal err: find max head failed")
	}
}

func (m *Migrate) quoteMergeDownRevisionAndEnd(str []string) string {
	var ret []string
	for _, s := range str {
		ret = append(ret, fmt.Sprintf("\n\t\t\t%v", m.quoteStrToMigrations(s)))
	}
	return fmt.Sprintf("[]string{%v,\n\t\t},\n\t}\n}", strings.Join(ret, ", "))
}

func (m *Migrate) MigrationsMergeContent(fn string, heads []string) string {
	return fmt.Sprintf("package migrations\n\nimport %v\n\nfunc (*Migrations) Migration_%v() *migrate.Operations {\n\treturn &migrate.Operations{\n\t\tOperations: []*migrate.Operation{},\n\t\tRevision:   %v,\n\t\tDownRevision: %v",
		m.quoteStrToMigrations(m.PackagePath), fn, m.quoteStrToMigrations(fn), m.quoteMergeDownRevisionAndEnd(heads))
}

func (m *Migrate) Merge() {
	defer m.handleErr()()

	node := m.GetOperationsTree(false)
	heads := m.HeadToString(node)
	if len(heads) < 2 {
		panic(fmt.Sprint("No Multiple heads need to merge"))
	}
	mh := maxHead(heads)
	fn := m.genMigrationFileName(mh)
	content := m.MigrationsMergeContent(fn, heads)
	m.write(content, fn)
}
