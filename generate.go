package goresminpack

import (
	"net/http"
	"path"
	"strings"

	"github.com/shurcooL/vfsgen"
)

const hotSwapCode = `
package {{ .package }}
+build debug,dev

// Locate current directory? Do I even need this? Yes
_, gofile, _, ok := runtime.Caller(1)
if ok {
	directory, _ = filepath.Abs(filepath.Join(filepath.Dir(gofile), filepath.FromSlash(directory)))
}
`

// Generate packs the resources into a Go file.
func Generate(assetDirectory, outputFilePath, outputObjectName string) (err error) {
	d := http.Dir(assetDirectory)
	return vfsgen.Generate(d, vfsgen.Options{
		Filename:     outputFilePath,
		PackageName:  strings.ToLower(path.Base(path.Dir(outputFilePath))),
		VariableName: outputObjectName,
		BuildTags:    "!debug,!dev",
	})
}
