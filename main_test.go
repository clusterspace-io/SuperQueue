package main

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("## Beginning Main Tests")
	exitVal := m.Run()
	fmt.Print("## Tests exited with status code", exitVal)
}

func HandleTestError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
