package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/pkg/lpm"
)

const DefaultVersionTemplate = "" +
	"{{with .Name}}{{printf \"%s \" .}}{{end}}" +
	"{{with .Short}}{{printf \"(%s) \" .}}{{end}}" +
	"{{printf \"version %s\" .Version}}\n"

func init() {
	//Global params
	//rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(
		&lpm.GetCliArgs().Category,
		"category",
		"",
		"extension, driver, or utility")

	rootCmd.Version = lpm.VersionNumber

	rootCmd.SetVersionTemplate(DefaultVersionTemplate)
}

var rootCmd = &cobra.Command{
	Use:   "lpm",
	Short: "Liquibase Package Manager",
	Long: `Easily manage external dependencies for Database Development.
Search for, install, and uninstall liquibase drivers, extensions, and utilities.`,
}

//Execute main entry point for CLI from root
func Execute(path string) error {
	ctx := lpm.NewContext(path)
	err := ctx.Initialize()
	if err != nil {
		err = fmt.Errorf("unable to initialize Context when executing root command; %w",
			err)
		goto end
	}
	err = rootCmd.ExecuteContext(ctx)
end:
	return err
}
