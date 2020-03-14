package goresminpack

import "testing"

func TestGenerator(t *testing.T) {
	err := Generate(`tests/pregeneration`, `tests/data.gen.go`, `Data`)
	if err != nil {
		panic(err)
	}
}
