package core

import "fmt"

func (m *Migrate) List() {
	defer m.handleErr()()
	for _, name := range m.AppliedOrdered() {
		fmt.Printf("Applied: %v\n", name)
	}
	for _, ops := range m.UnApplied() {
		fmt.Printf("Unapplied: %v\n", ops.Revision)
	}
}

func (m *Migrate) handleErr () func() {
	return func() {
		if x := recover(); x != nil {
			if v, ok := x.(error); ok {
				fmt.Println(v.Error())
			} else {
				fmt.Println(x)
			}
		}
	}
}