package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github/dkotik/zassets"
	"github/dkotik/zassets/compile"

	"github.com/spf13/cobra"
)

func main() {
	var CLI = &cobra.Command{
		Use:     `zassets`,
		Version: `0.0.1 Alpha`,
		Short:   `Compile and embed a resource directory into a Go module.`,
		Long: `Compile and embed a resource directory into a Go module.
Embedded resources are stored in an object
that satisfies http.FileSystem interface.`,
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := cmd.PersistentFlags().GetStringSlice(`embed`)
			if err != nil {
				cmd.Help()
				return
			}
			variable, err := cmd.PersistentFlags().GetString(`var`)
			if err != nil || variable == "" {
				fmt.Print("Error: variable name must be specified.\n")
				return
			}
			p, err := cmd.PersistentFlags().GetString(`package`)
			if err != nil || p == "" {
				fmt.Print("Error: package name must be specified.\n")
				return
			}
			if ok, _ := cmd.PersistentFlags().GetBool(`refine`); ok {
				t, err := ioutil.TempDir(os.TempDir(), `zassets-*`)
				if err != nil {
					fmt.Printf("Error: %s.\n", err.Error())
					return
				}
				c, _ := compile.NewCompiler()
				err = c.Run(t, entries...)
				defer os.RemoveAll(t)
				if err != nil {
					fmt.Printf("Error: %s.\n", err.Error())
					return
				}
				entries = []string{t}
			}
			comment, _ := cmd.PersistentFlags().GetString(`comment`)
			tags, _ := cmd.PersistentFlags().GetStringSlice(`tags`)
			err = zassets.EmbedAll(os.Stdout, &zassets.EmbedValues{
				Variable: variable,
				Package:  p,
				Comment:  comment,
				Tags:     strings.Join(tags, `,`)},
				entries...)
			if err != nil {
				fmt.Printf("Error: %s.\n", err.Error())
			}
		},
	}

	CLI.PersistentFlags().StringSliceP(`embed`, `e`, []string{}, `A list of files and directories containing common static assets.`)
	CLI.PersistentFlags().BoolP(`refine`, `r`, false, `Apply default refiners to assets before embedding.`)
	CLI.PersistentFlags().StringP(`var`, `v`, ``, `Assets will be accessible using this variable name.`)
	CLI.PersistentFlags().StringP(`package`, `p`, ``, `Assets will belong to this package.`)
	CLI.PersistentFlags().StringSliceP(`tags`, `t`, []string{}, `Specify build tags.`)
	CLI.PersistentFlags().StringP(`comment`, `c`, ``, `Include a comment.`)
	CLI.PersistentFlags().BoolP(`debug`, `d`, os.Getenv(`DEBUG`) != ``, `Make output as readable as possible.`)
	CLI.Execute()
}
