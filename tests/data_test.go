package tests

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func dump(p string) {
	f, err := Data.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(os.Stdout, f)
	fmt.Print("\n- - -\n")
	if err != nil {
		panic(err)
	}
}

func TestData(t *testing.T) {
	dump(`/test.css`)
	dump(`/template.tmpl`)
	dump(`/test.js`)
	t.Fail()
}
