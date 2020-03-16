package zassets

import (
	"os"
	"testing"
)

func TestEmbed(t *testing.T) {
	err := Embed(os.Stdout, &EmbedValues{
		Name:    "Assets",
		Package: "test",
		Comment: "comment\ncomment2\ncomment3",
		Data:    []byte("sdkfj slakdjflksdjf lsdk jflsdkjf lsadkfjsdlfjk"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
