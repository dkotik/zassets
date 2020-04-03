package compile

import "testing"

func TestCompiler(t *testing.T) {
	c, err := NewCompiler(WithDefaultOptions())
	if err != nil {
		t.Fatal(err)
	}

	i, err := NewIterator([]string{`../tests`, `text.go`}, []string{`\.go$`})
	err = c.Run(`/tmp/zassets`, i)
	// NewDebugger([]string{`../tests`, `text.go`}, []string{}, c)
	//
	// wait := make(chan bool)
	// <-wait // waiting

	if err != nil {
		t.Fatal(err)
	}
}
