package core

type MIGRATE struct {
	Command   string
	Reversion string
	DBDsn     string
}

func (m *Migrate) Run(command, reversion string) {
	switch command {
	case "migrate":
		m.Migrate()
	case "makemigrations":
		m.MakeMigrations()
	case "merge":
		m.Merge()
	case "fake":
		if reversion != "" {
			m.Fake(reversion)
		}
	case "list":
		m.List()
	case "downgrade":
		if reversion != "" {
			m.DownGrade(reversion)
		}
	}
}
