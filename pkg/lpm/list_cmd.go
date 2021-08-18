package lpm

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(listCmd)
	// @TODO make an CliArgs struct for `global` and other CLI args
	listCmd.Flags().BoolVarP(&cliArgs.Global, "global", "g", false, "list global packages")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Installed Packages",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		var files ClasspathFiles
		var cp string
		var err error
		var p Package
		var v Version
		var out string
		// Collect installed packages
		var installed Packages

		ctx := ContextFromCobraCommand(cmd)

		files, cp, err = ctx.GetClasspathFiles()
		if err != nil {
			ctx.Error("unable to get files for classpath '%s'; %w", cp, err)
			goto end
		}

		for _, p = range ctx.GetPackages() {
			v = p.GetInstalledVersion(files)
			if !files.VersionExists(v) {
				continue
			}
			installed = append(installed, p)
		}

		// Format output
		fmt.Println(cp)
		if len(installed) == 0 {
			os.Exit(1)
		}

		for _, out = range installed.Display(files) {
			fmt.Println(out)
		}
	end:
	},
}
