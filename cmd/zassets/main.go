package main

import (
	"os"

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

			cmd.Help()
		},
	}

	// CLI.PersistentFlags().StringP(`output`, `o`, ``, `Override default output file name.`)
	CLI.PersistentFlags().StringSliceP(`include`, `i`, []string{}, `A list of directories containing common static assets.`)
	CLI.PersistentFlags().BoolP(`debug`, `d`, os.Getenv(`DEBUG`) != ``, `Make output as readable as possible.`)
}
