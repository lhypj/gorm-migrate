package core

import "fmt"

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
		panic(fmt.Sprintf("version %v is blank", version))
	}
	if m.Applied()[version] == APPLIED {
		panic(fmt.Sprintf("version %v is applied", version))
	}
	if m.unAppliedList()[version] == 0 {
		panic(fmt.Sprintf("version %v not found", version))
	}
	if err := m.DB.Create(&OrmMigrations{Name: version}); err != nil {
		panic(err)
	}
	fmt.Printf("Fake: %v successful!", version)
}
