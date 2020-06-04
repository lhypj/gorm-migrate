package core

import (
	"flag"
)

type Args struct {
	Command   string `required:"true"`
	Reversion string
}

func (m *Migrate) Run() {
	var showUsage bool
	var command, reversion string
	flag.StringVar(&command, "c", "", "commands: makemigrations migrate merge fake list downgrade")
	flag.StringVar(&reversion, "v", "", "specify a reversion to downgrade or fake")

	flag.Parse()
	switch command {
	case "migrate":
		m.Migrate()
	case "makemigrations":
		m.MakeMigrations()
	case "merge":
		m.Merge()
	case "fake":
		if reversion == "" {
			showUsage = true
		} else {
			m.Fake(reversion)
		}
	case "list":
		m.List()
	case "downgrade":
		if reversion == "" {
			showUsage = true
		} else {
			m.DownGrade(reversion)
		}
	default:
		showUsage = true
	}
	if showUsage {
		flag.Usage()
	}
}
