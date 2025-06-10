package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
)

// ParseVersion parses a version string and returns a version object
// Handles formats like "0.2.9", "v0.2.9", and pre-releases
func ParseVersion(versionStr string) (*version.Version, error) {
	// Trim whitespace and normalize the version string
	versionStr = strings.TrimSpace(versionStr)
	
	// The go-version library handles "v" prefix automatically
	return version.NewVersion(versionStr)
}

// CompareVersions compares two version strings
// Returns:
//   -1 if current < latest
//    0 if current == latest
//    1 if current > latest
func CompareVersions(current, latest string) (int, error) {
	currentVer, err := ParseVersion(current)
	if err != nil {
		return 0, err
	}

	latestVer, err := ParseVersion(latest)
	if err != nil {
		return 0, err
	}

	return currentVer.Compare(latestVer), nil
}

// IsUpdateAvailable checks if an update is available
// Returns true if latest version is greater than current version
func IsUpdateAvailable(current, latest string) (bool, error) {
	result, err := CompareVersions(current, latest)
	if err != nil {
		return false, err
	}

	// Update is available if current < latest (result == -1)
	updateAvailable := result < 0
	
	// Only log in debug mode to avoid noise during normal operations
	if current == "dev" {
		// Development versions always consider updates available
		return true, nil
	}
	
	return updateAvailable, nil
}

// GetCurrentVersion returns the current lpm version from the VERSION file
// This function reads the VERSION file from the expected location
func GetCurrentVersion() string {
	// Try to find the VERSION file in the expected locations
	possiblePaths := []string{
		"internal/app/VERSION",
		"../VERSION",
		"../app/VERSION",
		filepath.Join(os.Getenv("GOPATH"), "src", "package-manager", "internal", "app", "VERSION"),
	}
	
	for _, path := range possiblePaths {
		if data, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	
	// If VERSION file cannot be found, return a fallback
	// In a real application, this should integrate with the build system
	return "dev"
}