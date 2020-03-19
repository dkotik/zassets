package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/dkotik/zassets"
	"github.com/dkotik/zassets/compile"

	"github.com/spf13/cobra"
)

func endOnError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s.\n", err.Error())
	// panic(err)
	os.Exit(1)
}

func main() {
	refine, embed, debug := false, false, false
	ignore := make([]string, 0)
	ev := &zassets.EmbedValues{}
	var CLI = &cobra.Command{
		Use:     `zassets`,
		Version: `0.0.1 Alpha`,
		Short:   `Compile and embed a resource directory into a Go module.`,
		Long: `Compile and embed a resource directory into a Go module.
Embedded resources are stored in an object
that satisfies http.FileSystem interface.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			iterator, err := compile.NewIterator(args, ignore)
			endOnError(err)

			if refine && !embed {
				c, err := compile.NewCompiler(compile.WithDefaultOptions())
				endOnError(err)
				if debug {
					compile.WithDebug()(c)
				}
				o, _ := cmd.PersistentFlags().GetString(`output`)
				err = c.Run(o, iterator)
				endOnError(err)
				return
			}

			if refine {
				t, err := ioutil.TempDir(os.TempDir(), `zassets-*`)
				endOnError(err)
				c, err := compile.NewCompiler(compile.WithDefaultOptions())
				endOnError(err)
				if debug {
					compile.WithDebug()(c)
				}
				err = c.Run(t, iterator)
				defer os.RemoveAll(t)
				endOnError(err)
				iterator = &compile.Iterator{
					Entries: []string{t},
					Ignore:  []*regexp.Regexp{}}
			}
			w := os.Stdout
			if o, err := cmd.PersistentFlags().GetString(`output`); o != "" {
				endOnError(err)
				w, err = os.Create(o)
				endOnError(err)
				defer w.Close()
			}
			endOnError(zassets.EmbedAll(w, ev, iterator))
		},
	}
	CLI.PersistentFlags().StringP(`output`, `o`, ``, `Write program output to this location.`)
	CLI.PersistentFlags().BoolVarP(&embed, `embed`, `e`, false, `Embed provided files and directories.`)
	CLI.PersistentFlags().BoolVarP(&refine, `refine`, `r`, false, `Apply default refiners to assets before embedding.`)
	CLI.PersistentFlags().StringVarP(&ev.Variable, `var`, `v`, ``, `Assets will be accessible using this variable name.`)
	CLI.PersistentFlags().StringVarP(&ev.Package, `package`, `p`, ``, `Assets will belong to this package.`)
	CLI.PersistentFlags().StringArrayVarP(&ev.Tags, `tags`, `t`, []string{}, `Specify build tags.`)
	CLI.PersistentFlags().StringVarP(&ev.HashAlgorythm, `hashwith`, `hw`, ``, `Include a hash table in the embedded output. Choose from xx, md5, and sha256.`)
	CLI.PersistentFlags().StringArrayVarP(&ignore, `ignore`, `i`, []string{}, `Skip files and directories that match provided patterns.`)
	CLI.PersistentFlags().StringVarP(&ev.Comment, `comment`, `c`, ``, `Include a comment.`)
	CLI.PersistentFlags().BoolVarP(&debug, `debug`, `d`, os.Getenv(`DEBUG`) != ``, `Write the contents of refined files as readable as possible.`)
	CLI.Execute()
}
