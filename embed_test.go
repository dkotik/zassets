package zassets

import (
	"os"
	"testing"
)

func TestEmbedAll(t *testing.T) {
	err := EmbedAll(os.Stdout, &EmbedValues{
		Variable: "Assets",
		Package:  "tests",
		Comment:  "comment\ncomment2\ncomment3",
		Tags:     []string{`dev`, `debug`},
	}, `tests/go.mod`, `tests/data`)
	if err != nil {
		t.Fatal(err)
	}
}
