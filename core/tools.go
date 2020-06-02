package core

import "fmt"

func (m *Migrate) List() {
	if rows, err := m.DB.Model(&OrmMigrations{}).Select("name").Order("id").Rows(); err != nil {
		panic(err)
	} else {
		var name string
		for rows.Next() {
			if err := rows.Scan(&name); err != nil {
				panic(err)
			} else {
				fmt.Printf("Applied: %v\n", name)
			}
		}
	}
	for _, ops := range m.UnApplied() {
		fmt.Printf("Unapplied: %v\n", ops.Revision)
	}
}
