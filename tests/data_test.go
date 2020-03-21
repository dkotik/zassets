package tests

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/dkotik/zassets"
)

func dump(p string) {
	f, err := Assets.Open(p)
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
	zassets.Walk(Assets, `/`, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dump(path)
		}
		return nil
	}))
	// dump(`/test.css`)
	// dump(`/template.tmpl`)
	// dump(`/test.js`)
	// t.Fail()
}
