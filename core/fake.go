package core

import (
	"fmt"
)

func (m *Migrate) unAppliedList() map[string]int {
	ret := make(map[string]int)
	for _, op := range m.UnApplied() {
		if op.Revision != ""{
			ret[op.Revision] = 1
		}
	}
	return ret
}

func (m *Migrate) Fake(version string) {
	if version == "" {
		fmt.Printf("version %v is blank", version)
		return
	}
	if m.Applied()[version] == APPLIED {
		fmt.Printf("version %v is applied", version)
		return
	}
	if m.unAppliedList()[version] == 0 {
		fmt.Printf("version %v not found", version)
		return
	}
	if err := m.DB.Create(&OrmMigrations{Name: version}).Error; err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Fake: %v successful!", version)
	}
}
