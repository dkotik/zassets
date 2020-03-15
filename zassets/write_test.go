package zassets

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestZipWrite(t *testing.T) {
	w, err := os.Create(`tests/test.zip`)
	defer w.Close()
	err = write(w, http.Dir(`../tests`))
	if err != nil {
		panic(err)
	}
	w.Close()

	s := Must(FromArchive(`tests/test.zip`))
	// s.Open(`test`)
	Walk(s, `/`, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		log.Println("detected file:", path, info.IsDir(), info.ModTime())
		return nil
	}))
}
