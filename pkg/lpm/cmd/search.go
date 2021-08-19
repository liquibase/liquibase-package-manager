package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/pkg/lpm"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

// searchCmd represents the `search` command
var searchCmd = &cobra.Command{
	Use:   "search [PACKAGE]",
	Short: "Search for Packages",

	Run: func(cmd *cobra.Command, args []string) {
		var name string
		var files lpm.ClasspathFiles
		var cp string
		var err error
		var found lpm.Packages

		ctx := lpm.GetContextFromCommand(cmd)

		files, cp, err = ctx.GetClasspathFiles()
		if err != nil {
			ctx.Error("unable to get files for classpath '%s'; %w", cp, err)
			goto end
		}

		if len(args) > 0 {
			name = args[0]
		}

		if len(name) < 3 {
			fmt.Println("Minimum of 3 characters required for search.")
			goto end
		}

		found = ctx.SearchPackages(name)

		if len(found) == 0 {
			fmt.Println("No results found.")
			goto end
		}

		found.Println(files)

	end:
	},
}
