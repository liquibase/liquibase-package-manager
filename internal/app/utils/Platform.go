package utils

import (
	"fmt"
	"os"
	"runtime"
)

// PlatformUtil struct for platform detection utilities
type PlatformUtil struct{}

// GetCurrentPlatform returns the current OS/arch combination using the naming convention
// that matches the release asset names
func (p PlatformUtil) GetCurrentPlatform() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "darwin":
		switch goarch {
		case "amd64":
			return "darwin", nil
		case "arm64":
			return "darwin-arm64", nil
		default:
			return "", fmt.Errorf("unsupported darwin architecture: %s", goarch)
		}
	case "linux":
		switch goarch {
		case "amd64":
			return "linux", nil
		case "arm64":
			return "linux-arm64", nil
		case "s390x":
			return "s390x", nil
		default:
			return "", fmt.Errorf("unsupported linux architecture: %s", goarch)
		}
	case "windows":
		switch goarch {
		case "amd64":
			return "windows", nil
		default:
			return "", fmt.Errorf("unsupported windows architecture: %s", goarch)
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}

// GetBinaryName returns the binary name for the current platform
// Returns "lpm.exe" for Windows, "lpm" for all other platforms
func (p PlatformUtil) GetBinaryName() string {
	if runtime.GOOS == "windows" {
		return "lpm.exe"
	}
	return "lpm"
}

// GetReleaseAssetName returns the GitHub release asset name for the current platform
// using the naming convention: lpm-{version}-{platform}.zip
func (p PlatformUtil) GetReleaseAssetName(version string) string {
	platform, err := p.GetCurrentPlatform()
	if err != nil {
		// Fallback to a generic name if platform detection fails
		return fmt.Sprintf("lpm-%s.zip", version)
	}
	return fmt.Sprintf("lpm-%s-%s.zip", version, platform)
}

// GetExecutablePath returns the path to the current running binary
func (p PlatformUtil) GetExecutablePath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return executable, nil
}