package core

import "fmt"

func (m *Migrate) Run(command, reversion string) {
	switch command {
	case "migrate":
		m.Migrate()
	case "makemigrations":
		m.MakeMigrations()
	case "merge":
		m.Merge()
	case "fake":
		m.Fake(reversion)
	case "list":
		m.List()
	case "downgrade":
		m.DownGrade(reversion)
	default:
		fmt.Println("Command migrate、makemigrations、merge、fake、list、downgrade")
	}
}
