package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/spf13/cobra"
	"package-manager/internal/app"
	"package-manager/internal/app/dependencies"
	"package-manager/internal/app/errors"
	"package-manager/internal/app/packages"
	"package-manager/internal/app/utils"
)

// upgradeCmd represents the update command
var upgradeCmd = &cobra.Command{
	Use:     "upgrade [PACKAGE]...",
	Short:   "Upgrades lpm itself or installed packages to the latest versions",
	Aliases: []string{"up"},
	Long: `Upgrade lpm itself or installed packages:

  lpm upgrade                    # Update lpm itself to latest version
  lpm upgrade <package>          # Update specific package(s)
  lpm upgrade --version=X.Y.Z    # Update lpm to specific version
  lpm upgrade --check            # Check for updates without installing
  lpm upgrade --all              # Update both lpm and all packages

Debug logging: This command outputs detailed [DEBUG] and [ERROR] messages 
to help with troubleshooting. These include version comparisons, GitHub API 
calls, download progress, and update operations.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Handle different upgrade scenarios based on arguments and flags
		if len(args) == 0 && !upgradeAll {
			// No arguments - upgrade lpm itself
			handleSelfUpgrade(cmd)
			return
		}

		if upgradeAll {
			// Upgrade both lpm and all packages
			handleSelfUpgrade(cmd)
			fmt.Println()
		}

		// Handle package upgrades (existing functionality or --all flag)
		if len(args) > 0 || upgradeAll {
			var outdated packages.Packages
			
			if len(args) > 0 {
				// Specific packages requested - check each one
				for _, packageName := range args {
					p := packs.GetByName(packageName)
					if p.Name == "" {
						errors.Exit("Package '"+packageName+"' not found.", 1)
					}
					if !p.InClassPath(app.ClasspathFiles) {
						errors.Exit("Package '"+packageName+"' is not installed.", 1)
					}
					installed := p.GetInstalledVersion(app.ClasspathFiles)
					latest := p.GetLatestVersion(liquibase.Version)
					if latest.Tag != "" && installed.Tag != latest.Tag {
						outdated = append(outdated, p)
					}
				}
			} else {
				// --all flag - get all outdated packages
				outdated = packs.GetOutdated(liquibase.Version, app.ClasspathFiles)
			}
			
			if len(outdated) == 0 {
				if len(args) > 0 {
					fmt.Println("All specified packages are already up to date.")
				} else {
					fmt.Println("You have no outdated packages installed.")
					if !upgradeAll {
						fmt.Println(app.Classpath)
						os.Exit(0)
					}
				}
				return
			}
		var r []string
		var prefix string
		r = append(r, fmt.Sprintf("%-4s %-38s %-38s %s", "   ", "Package", "Installed", "Available"))
		for i, p := range outdated {
			ins := p.GetInstalledVersion(app.ClasspathFiles)
			latest := p.GetLatestVersion(liquibase.Version)
			if (i + 1) == len(outdated) {
				prefix = "└──"
			} else {
				prefix = "├──"
			}
			r = append(r, fmt.Sprintf("%-4s %-38s %-38s %s", prefix, p.Name, ins.Tag, latest.Tag))
		}
		fmt.Println("You have " + strconv.Itoa(len(outdated)) + " outdated package(s) installed.")
		fmt.Println(app.Classpath)
		for _, out := range r {
			fmt.Println(out)
		}
		if !dryRun {
			d := dependencies.Dependencies{}
			if !global {
				d.Read()
			}
			for _, p := range outdated {
				ins := p.GetInstalledVersion(app.ClasspathFiles)
				latest := p.GetLatestVersion(liquibase.Version)
				fmt.Println()
				fmt.Println("removing " + p.Name + "@" + ins.Tag + " from classpath")
				err := p.Remove(app.Classpath, ins)
				if err != nil {
					errors.Exit("Unable to remove "+ins.GetFilename()+" from classpath.", 1)
				}
				fmt.Println(ins.GetFilename() + " successfully uninstalled from classpath.")
				if !global {
					d.Remove(p.Name)
				}
				fmt.Println()
				fmt.Println("adding " + p.Name + "@" + latest.Tag + " to classpath")
				if !latest.PathIsHTTP() {
					latest.CopyToClassPath(app.Classpath)
				} else {
					latest.DownloadToClassPath(app.Classpath)
				}
				fmt.Println(latest.GetFilename() + " successfully installed in classpath.")
				d.Dependencies = append(d.Dependencies, dependencies.Dependency{p.Name: latest.Tag})
			}
			if !global {
				d.Write()
			}
		}
	}
	},
}

// Package-level variables for new flags
var (
	upgradeVersion string
	upgradeCheck   bool
	upgradeAll     bool
)

// handleSelfUpgrade handles upgrading lpm itself
func handleSelfUpgrade(cmd *cobra.Command) {
	currentVersion := app.Version()
	if currentVersion == "" {
		currentVersion = "dev"
	}

	fmt.Printf("Current lpm version: %s\n", currentVersion)

	// Initialize utilities
	githubUtil := utils.NewGitHubUtil()
	versionUtil := &utils.PlatformUtil{}
	downloadUtil := utils.NewDownloadUtil()
	updateUtil := utils.NewUpdateUtil()

	// Repository information for lpm
	owner := "liquibase"
	repo := "liquibase-package-manager"

	var targetVersion string
	var release *github.RepositoryRelease
	var err error

	if upgradeVersion != "" {
		// Upgrade to specific version
		targetVersion = upgradeVersion
		fmt.Printf("Checking for lpm version %s...\n", targetVersion)
		fmt.Printf("[DEBUG] Fetching release %s from repository %s/%s\n", targetVersion, owner, repo)
		release, err = githubUtil.GetRelease(owner, repo, targetVersion)
		if err != nil {
			fmt.Printf("[ERROR] GitHub API request failed for version %s: %s\n", targetVersion, err.Error())
			errors.Exit(fmt.Sprintf("Failed to find version %s: %s", targetVersion, err.Error()), 1)
		}
		fmt.Printf("[DEBUG] Successfully found release %s\n", targetVersion)
	} else {
		// Upgrade to latest version
		fmt.Println("Checking for latest lpm version...")
		fmt.Printf("[DEBUG] Fetching latest release from repository %s/%s\n", owner, repo)
		release, err = githubUtil.GetLatestRelease(owner, repo)
		if err != nil {
			fmt.Printf("[ERROR] GitHub API request failed for latest release: %s\n", err.Error())
			errors.Exit("Failed to check for latest version: "+err.Error(), 1)
		}
		targetVersion = strings.TrimPrefix(release.GetTagName(), "v")
		fmt.Printf("[DEBUG] Latest release found: %s\n", targetVersion)
	}

	// Check if update is needed
	if currentVersion != "dev" {
		fmt.Printf("[DEBUG] Comparing versions: current=%s, target=%s\n", currentVersion, targetVersion)
		isUpdateAvailable, err := utils.IsUpdateAvailable(currentVersion, targetVersion)
		if err != nil {
			fmt.Printf("[ERROR] Version comparison failed: %s\n", err.Error())
			errors.Exit("Failed to compare versions: "+err.Error(), 1)
		}
		if !isUpdateAvailable {
			fmt.Printf("lpm is already up to date (version %s)\n", currentVersion)
			fmt.Printf("[DEBUG] No update needed: current version %s >= target version %s\n", currentVersion, targetVersion)
			return
		}
		fmt.Printf("[DEBUG] Update available: %s -> %s\n", currentVersion, targetVersion)
	}

	fmt.Printf("Update available: %s -> %s\n", currentVersion, targetVersion)

	// If check flag is set, just report and exit
	if upgradeCheck {
		fmt.Println("Use 'lpm upgrade' to install the update.")
		return
	}

	// Proceed with download and update
	if dryRun {
		fmt.Printf("[DRY RUN] Would upgrade lpm from %s to %s\n", currentVersion, targetVersion)
		return
	}

	// Get platform information
	fmt.Printf("[DEBUG] Detecting platform information\n")
	platform, err := versionUtil.GetCurrentPlatform()
	if err != nil {
		fmt.Printf("[ERROR] Platform detection failed: %s\n", err.Error())
		errors.Exit("Failed to detect platform: "+err.Error(), 1)
	}
	fmt.Printf("[DEBUG] Platform detected: %s\n", platform)

	fmt.Printf("Downloading lpm %s for %s...\n", targetVersion, platform)

	// Download the binary
	fmt.Printf("[DEBUG] Starting download for release %s\n", release.GetTagName())
	newBinaryPath, err := downloadUtil.DownloadBinaryUpdate(release, platform)
	if err != nil {
		fmt.Printf("[ERROR] Download failed: %s\n", err.Error())
		errors.Exit("Failed to download update: "+err.Error(), 1)
	}
	defer os.Remove(newBinaryPath) // Clean up downloaded file
	fmt.Printf("[DEBUG] Download completed, binary saved to: %s\n", newBinaryPath)

	fmt.Println("Installing update...")

	// Perform atomic update
	fmt.Printf("[DEBUG] Starting atomic update process\n")
	err = updateUtil.PerformAtomicUpdate(newBinaryPath)
	if err != nil {
		// Special handling for Windows where process restart is needed
		if strings.Contains(err.Error(), "Windows update requires process restart") {
			fmt.Printf("[DEBUG] Windows update staged successfully, requiring restart\n")
			fmt.Println("\nUpdate staged successfully!")
			fmt.Println("Please restart lpm to complete the update process.")
			os.Exit(0)
		}
		fmt.Printf("[ERROR] Update installation failed: %s\n", err.Error())
		errors.Exit("Failed to install update: "+err.Error(), 1)
	}

	fmt.Printf("[DEBUG] Update completed successfully\n")
	fmt.Printf("\nlpm successfully updated to version %s!\n", targetVersion)
	fmt.Println("Run 'lpm --version' to verify the update.")
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVarP(&global, "global", "g", false, "upgrade global packages")
	upgradeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "output changes without applying")
	upgradeCmd.Flags().StringVar(&upgradeVersion, "version", "", "upgrade lpm to specific version")
	upgradeCmd.Flags().BoolVar(&upgradeCheck, "check", false, "check for updates without installing")
	upgradeCmd.Flags().BoolVar(&upgradeAll, "all", false, "upgrade both lpm and all packages")
}
