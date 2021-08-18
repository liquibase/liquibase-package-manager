package lpm

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(
		&cliArgs.Global,
		"global",
		"g",
		false,
		"add package globally")
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [PACKAGE]...",
	Short: "Add Packages",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var p Package
		var v Version
		var err error

		ctx := ContextFromCobraCommand(cmd)

		dd := NewDependencies()

		if ctx.FileSource == LocalFiles {
			err = dd.ReadManifest(ctx)
		}
		if err != nil {
			ctx.Error("unable to read local classpath files")
			goto end
		}

		for _, name := range args {

			err = maybeAddPackage(ctx, name)
			if err != nil {
				ctx.Error("unable to add package %s; %w", name, err)
				continue
			}

			dd.Append(NewDependency(p.Name, v.Tag))

			fmt.Printf("%s successfully installed in classpath.\n",
				v.GetFilename())

		}

		if ctx.FileSource == LocalFiles {
			//Add package to local manifest
			if !dd.FileExists(ctx) {
				err = dd.CreateManifestFile(ctx)
			}
			if err != nil {
				ctx.Error("unable to create local manifest %s", ctx.GetManifestFilepath())
				goto end
			}
			err = dd.WriteManifest(ctx)
			if err != nil {
				ctx.Error("unable to write to local manifest %s", ctx.GetManifestFilepath())
				goto end
			}

		}

		// Output helper for JAVA_OPTS
		ctx.PrintJavaOptsHelper()
	end:
	},
}

func maybeAddPackage(ctx *Context, name string) error {
	var cp string
	var msg string
	var files ClasspathFiles

	p, v, err := ctx.GetPackageAndVersion(name)

	if p.Name == "" {
		msg = fmt.Sprintf("package '%s' not found", name)
		goto end
	}

	files, cp, err = ctx.GetClasspathFiles()
	if err != nil {
		msg = fmt.Sprintf("unable to get classpath files")
		goto end
	}

	if files.VersionExists(v) {
		msg = fmt.Sprintf("%s is already installed", name)
		goto end
	}

	err = v.CopyFilesToClassPath(cp)
	if err != nil {
		msg = fmt.Sprintf("unable to copy %s file to classpath", name)
		goto end
	}

end:
	if msg != "" {
		err = fmt.Errorf("%s when attempting to add", msg)
	}
	return err
}
