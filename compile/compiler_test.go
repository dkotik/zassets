package compile

import "testing"

func TestCompiler(t *testing.T) {
	c, err := NewCompiler(WithDebug())
	// c, err := NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Run(`/tmp/zassets`, `../tests/data`)
	if err != nil {
		t.Fatal(err)
	}
}
