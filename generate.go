package goresminpack

import (
	"fmt"
	"log"
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

// Generate minifies and packs the resources into the output file.
func Generate(assetDirectory, outputFilePath, outputObjectName string) (err error) {
	d := Dir(assetDirectory)
	// d := http.Dir(assetDirectory)
	err = vfsgen.Generate(d, vfsgen.Options{
		Filename:        outputFilePath,
		PackageName:     strings.ToLower(path.Base(path.Dir(outputFilePath))),
		VariableName:    outputObjectName,
		VariableComment: fmt.Sprintf("%s holds", outputObjectName),
		BuildTags:       "!debug,!dev",
	})
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
