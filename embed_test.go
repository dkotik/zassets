package zassets

import (
	"os"
	"testing"

	"github.com/dkotik/zassets/compile"
)

func TestEmbedAll(t *testing.T) {
	i, _ := compile.NewIterator(
		[]string{`tests/go.mod`, `tests/data`},
		[]string{},
	)
	em := &Embed{
		Variable:      "Assets",
		Package:       "tests",
		Comment:       "comment\ncomment2\ncomment3",
		Tags:          []string{`dev`, `debug`},
		HashAlgorythm: `xx`,
	}
	em.SetTemplate("")
	err := em.Iterator(os.Stdout, i)
	if err != nil {
		t.Fatal(err)
	}
}
