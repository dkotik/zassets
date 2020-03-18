package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github/dkotik/zassets"
	"github/dkotik/zassets/compile"

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

// TODO: three execution paths: embed, compile, compile and embed

func main() {
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
			// entries, err := cmd.PersistentFlags().GetStringArray(`embed`)
			if len(args) == 0 {
				cmd.Help()
				return
			}
			if ok, _ := cmd.PersistentFlags().GetBool(`refine`); ok {
				t, err := ioutil.TempDir(os.TempDir(), `zassets-*`)
				endOnError(err)
				c, err := compile.NewCompiler(
					compile.WithIgnore(ignore...),
					compile.WithDefaultOptions())
				endOnError(err)
				if ok, _ = cmd.PersistentFlags().GetBool(`debug`); ok {
					compile.WithDebug()(c)
				}
				err = c.Run(t, args...)
				defer os.RemoveAll(t)
				endOnError(err)
				args = []string{t}
			}
			w := os.Stdout
			if o, err := cmd.PersistentFlags().GetString(`output`); o != "" {
				endOnError(err)
				w, err = os.Create(o)
				endOnError(err)
				defer w.Close()
			}
			err := zassets.EmbedAll(w, ev, args...)
			endOnError(err)
		},
	}
	CLI.PersistentFlags().StringP(`output`, `o`, ``, `Write program output to this location.`)
	CLI.PersistentFlags().BoolP(`embed`, `e`, false, `Embed provided files and directories.`)
	CLI.PersistentFlags().BoolP(`refine`, `r`, false, `Apply default refiners to assets before embedding.`)
	CLI.PersistentFlags().StringVarP(&ev.Variable, `var`, `v`, ``, `Assets will be accessible using this variable name.`)
	CLI.PersistentFlags().StringVarP(&ev.Package, `package`, `p`, ``, `Assets will belong to this package.`)
	CLI.PersistentFlags().StringArrayVarP(&ev.Tags, `tags`, `t`, []string{}, `Specify build tags.`)
	CLI.PersistentFlags().StringArrayVarP(&ignore, `ignore`, `i`, []string{}, `Skip files and directories that match provided patterns.`)
	CLI.PersistentFlags().StringVarP(&ev.Comment, `comment`, `c`, ``, `Include a comment.`)
	CLI.PersistentFlags().BoolP(`debug`, `d`, os.Getenv(`DEBUG`) != ``, `Write the contents of refined files as readable as possible.`)
	CLI.Execute()
}
