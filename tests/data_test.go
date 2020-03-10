package tests

import (
	"io"
	"os"
	"testing"
)

func TestData(t *testing.T) {
	f, err := Data.Open(`/test.sass.bahaha`)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		panic(err)
	}
	t.Fail()
}
