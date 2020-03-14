package goresminpack

import "testing"

func TestCompiler(t *testing.T) {
	// err := Generate(`tests/data`, `tests/data.gen.go`, `Data`)
	c := &Compiler{}
	err := c.Process(`tests/data`, `tests/pregeneration`)
	if err != nil {
		panic(err)
	}
}
