package lpm

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(
		&cliArgs.Global,
		"global",
		"g",
		false,
		"remove package globally")
}

// removeCmd represents the `remove` command
var removeCmd = &cobra.Command{
	Use:     "remove [PACKAGE]...",
	Short:   "Removes Package",
	Aliases: []string{"rm"},
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var files ClasspathFiles
		var cp string
		var p Package
		var v Version
		var err error

		ctx := ContextFromCobraCommand(cmd)

		d := NewDependencies()

		if ctx.FileSource != GlobalFiles {
			err = d.ReadManifest(ctx)
		}
		if err != nil {
			goto end
		}

		files, cp, err = ctx.GetClasspathFiles()
		if err != nil {
			ctx.Error("unable to get files for classpath '%s'; %w", cp, err)
			goto end
		}

		// Remove Each Package
		for _, name := range args {
			p = ctx.GetPackageByName(name)
			v = p.GetInstalledVersion(files)
			if p.Name == "" {
				ctx.Error("package '%s' not found.", name)
			}
			if !v.InClassPath(files) {
				ctx.Error("%s is not installed.", name)
			}
			err := os.Remove(cp + v.GetFilename())
			if err != nil {
				ctx.Error("unable to remove %s from classpath %s", v.GetFilename(), cp)
				continue
			}
			fmt.Printf("%s successfully uninstalled from classpath.\n", v.GetFilename())
			if ctx.FileSource != GlobalFiles {
				d.Remove(p.Name)
			}
		}

		if ctx.FileSource != GlobalFiles {
			err = d.WriteManifest(ctx)
		}
		if err != nil {
			ctx.Error("unable to write manifest %s for classpath %s",
				v.GetFilename(),
				cp)
			goto end
		}

	end:
		return

	},
}
