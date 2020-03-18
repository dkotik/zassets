package compile

import "testing"

func TestCompiler(t *testing.T) {
	c, err := NewCompiler(WithDebug(), WithDefaultOptions())
	// c, err := NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	i, err := NewIterator([]string{`../tests`, `text.go`}, []string{`\.go$`})
	err = c.Run(`/tmp/zassets`, i)
	if err != nil {
		t.Fatal(err)
	}
}
