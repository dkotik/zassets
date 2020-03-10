package goresminpack

import "testing"

func TestFileSystem(t *testing.T) {
	err := Generate(`tests/data`, `tests/data.gen.go`, `Data`)
	if err != nil {
		panic(err)
	}
}
