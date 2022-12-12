package commands

import (
	"fmt"
    "sort"
    "github.com/hashicorp/go-version"
    "package-manager/internal/app"
    "package-manager/internal/app/errors"
    "github.com/spf13/cobra"
)

// dedupeCmd represents the dedupe command
var dedupeCmd = &cobra.Command{
	Use:   "dedupe",
    Short: "Deduplicate Packages",
    Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

        for _, p := range packs.GetInstalled(app.ClasspathFiles) {
            var installed []*version.Version
            for _, v := range p.Versions {
                if v.InClassPath(app.ClasspathFiles) {
                    ver, _ := version.NewVersion(v.Tag)
                    installed = append(installed, ver)
                }
            }

            if len(installed) == 1 {
                continue
            }

            sort.Sort(sort.Reverse(version.Collection(installed)))
            fmt.Println(app.Classpath)
            var r []string

            var prefix string
            r = append(r, fmt.Sprintf("%-4s %-38s %s", "   ", "Package", "Installed"))
            for i, v := range installed {
                if (i + 1) == len(installed) {
                    prefix = "└──"
                } else {
                    prefix = "├──"
                }
                ver := p.GetVersion(v.Original())
                r = append(r, fmt.Sprintf("%-4s %-38s %s", prefix, p.Name, ver.Tag))
            }
            for _, out := range r {
                fmt.Println(out)
            }
            if !dryRun {
                for i, v := range installed {
                    if i == 0 {
                        continue
                    }
                    fmt.Println()
                    ver := p.GetVersion(v.Original())
                    fmt.Println("removing " + p.Name + "@" + ver.Tag + " from classpath")
                    err := p.Remove(app.Classpath, ver)
                    if err != nil {
                        errors.Exit("Unable to remove "+ver.GetFilename()+" from classpath.", 1)
                    }
                    fmt.Println(ver.GetFilename() + " successfully uninstalled from classpath.")
                }
            }
            fmt.Println()
        }
	},
}

func init() {
    rootCmd.AddCommand(dedupeCmd)
    dedupeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "output changes without applying")
}
