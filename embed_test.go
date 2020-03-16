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
	}, `tests/go.mod`, `tests/data`)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestEmbed(t *testing.T) {
// 	r, err := os.Open(`tests/test.zip`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer r.Close()
//
// 	w, err := os.Create(`tests/data.gen.go`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer w.Close()
//
// 	err = Embed(w, r, &EmbedValues{
// 		Name:    "Assets",
// 		Package: "tests",
// 		Comment: "comment\ncomment2\ncomment3",
// 	}, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
