package core

import (
	"fmt"
	"os"
)

func (m *Migrate)Run() {
	fmt.Println(os.Args)
}
