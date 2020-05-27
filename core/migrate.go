package core

import "fmt"

func (*Migrate) UnApplied() *[]string {
	var ret *[]string
	return ret
}

func (*Migrate) Fake() {

}

func (*Migrate) Migrate() {
	fmt.Println("done")
}
