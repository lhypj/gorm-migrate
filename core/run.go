package core

type M struct {
	Command   string `required:"true"`
	Reversion string
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
