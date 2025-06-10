package utils

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetCurrentPlatform(t *testing.T) {
	// Store original runtime values
	originalGOOS := runtime.GOOS
	originalGOARCH := runtime.GOARCH

	// Since we can't easily mock runtime.GOOS and runtime.GOARCH,
	// we'll test the actual current platform and verify it returns
	// a known valid platform string
	platformUtil := PlatformUtil{}
	platform, err := platformUtil.GetCurrentPlatform()

	if err != nil {
		t.Errorf("GetCurrentPlatform() returned error: %v", err)
	}

	// Verify the platform is one of the supported ones
	validPlatforms := []string{
		"darwin", "darwin-arm64",
		"linux", "linux-arm64", "s390x",
		"windows",
	}

	found := false
	for _, validPlatform := range validPlatforms {
		if platform == validPlatform {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("GetCurrentPlatform() returned unsupported platform: %s", platform)
	}

	// Test the mapping logic based on actual runtime values
	expectedPlatform := ""
	switch originalGOOS {
	case "darwin":
		if originalGOARCH == "amd64" {
			expectedPlatform = "darwin"
		} else if originalGOARCH == "arm64" {
			expectedPlatform = "darwin-arm64"
		}
	case "linux":
		if originalGOARCH == "amd64" {
			expectedPlatform = "linux"
		} else if originalGOARCH == "arm64" {
			expectedPlatform = "linux-arm64"
		} else if originalGOARCH == "s390x" {
			expectedPlatform = "s390x"
		}
	case "windows":
		if originalGOARCH == "amd64" {
			expectedPlatform = "windows"
		}
	}

	if expectedPlatform != "" && platform != expectedPlatform {
		t.Errorf("GetCurrentPlatform() = %s, want %s", platform, expectedPlatform)
	}
}

func TestGetBinaryName(t *testing.T) {
	platformUtil := PlatformUtil{}
	binaryName := platformUtil.GetBinaryName()

	// Test based on actual runtime OS
	if runtime.GOOS == "windows" {
		if binaryName != "lpm.exe" {
			t.Errorf("GetBinaryName() on Windows = %s, want lpm.exe", binaryName)
		}
	} else {
		if binaryName != "lpm" {
			t.Errorf("GetBinaryName() on non-Windows = %s, want lpm", binaryName)
		}
	}
}

func TestGetReleaseAssetName(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "simple version",
			version: "1.0.0",
		},
		{
			name:    "version with patch",
			version: "1.2.3",
		},
		{
			name:    "pre-release version",
			version: "1.0.0-rc.1",
		},
		{
			name:    "version with build metadata",
			version: "1.0.0+build.123",
		},
	}

	platformUtil := PlatformUtil{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assetName := platformUtil.GetReleaseAssetName(tt.version)

			// Verify the asset name follows the expected format
			expectedPrefix := "lpm-" + tt.version + "-"
			expectedSuffix := ".zip"

			if !strings.HasPrefix(assetName, expectedPrefix) {
				t.Errorf("GetReleaseAssetName() = %s, should start with %s", assetName, expectedPrefix)
			}

			if !strings.HasSuffix(assetName, expectedSuffix) {
				t.Errorf("GetReleaseAssetName() = %s, should end with %s", assetName, expectedSuffix)
			}

			// Extract platform from asset name
			platformPart := strings.TrimPrefix(assetName, expectedPrefix)
			platformPart = strings.TrimSuffix(platformPart, expectedSuffix)

			// Verify the platform part is valid
			validPlatforms := []string{
				"darwin", "darwin-arm64",
				"linux", "linux-arm64", "s390x",
				"windows",
			}

			found := false
			for _, validPlatform := range validPlatforms {
				if platformPart == validPlatform {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("GetReleaseAssetName() contains invalid platform: %s", platformPart)
			}
		})
	}
}

func TestGetExecutablePath(t *testing.T) {
	platformUtil := PlatformUtil{}
	execPath, err := platformUtil.GetExecutablePath()

	if err != nil {
		t.Errorf("GetExecutablePath() returned error: %v", err)
	}

	if execPath == "" {
		t.Error("GetExecutablePath() returned empty string")
	}

	// Verify the path exists
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		t.Errorf("GetExecutablePath() returned non-existent path: %s", execPath)
	}

	// Verify it's an absolute path
	if !strings.HasPrefix(execPath, "/") && !strings.Contains(execPath, ":\\") {
		t.Errorf("GetExecutablePath() should return absolute path, got: %s", execPath)
	}
}

// TestPlatformMappings tests the platform mapping logic more thoroughly
func TestPlatformMappings(t *testing.T) {
	tests := []struct {
		name        string
		goos        string
		goarch      string
		expected    string
		expectError bool
	}{
		// Darwin tests
		{
			name:     "darwin amd64",
			goos:     "darwin",
			goarch:   "amd64",
			expected: "darwin",
		},
		{
			name:     "darwin arm64",
			goos:     "darwin",
			goarch:   "arm64",
			expected: "darwin-arm64",
		},
		{
			name:        "darwin unsupported arch",
			goos:        "darwin",
			goarch:      "386",
			expectError: true,
		},
		// Linux tests
		{
			name:     "linux amd64",
			goos:     "linux",
			goarch:   "amd64",
			expected: "linux",
		},
		{
			name:     "linux arm64",
			goos:     "linux",
			goarch:   "arm64",
			expected: "linux-arm64",
		},
		{
			name:     "linux s390x",
			goos:     "linux",
			goarch:   "s390x",
			expected: "s390x",
		},
		{
			name:        "linux unsupported arch",
			goos:        "linux",
			goarch:      "arm",
			expectError: true,
		},
		// Windows tests
		{
			name:     "windows amd64",
			goos:     "windows",
			goarch:   "amd64",
			expected: "windows",
		},
		{
			name:        "windows unsupported arch",
			goos:        "windows",
			goarch:      "386",
			expectError: true,
		},
		// Unsupported OS
		{
			name:        "unsupported os",
			goos:        "freebsd",
			goarch:      "amd64",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily mock runtime.GOOS and runtime.GOARCH,
			// so we'll test the logic by examining what the function
			// would return for the current platform and verify it
			// matches our expected mappings

			// This is a conceptual test - in a real implementation,
			// you might want to refactor the function to accept
			// goos and goarch as parameters for easier testing
			platformUtil := PlatformUtil{}
			
			// For now, we'll just verify the current platform works
			if runtime.GOOS == tt.goos && runtime.GOARCH == tt.goarch {
				platform, err := platformUtil.GetCurrentPlatform()
				
				if tt.expectError {
					if err == nil {
						t.Errorf("Expected error for %s/%s, but got none", tt.goos, tt.goarch)
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error for %s/%s: %v", tt.goos, tt.goarch, err)
					}
					if platform != tt.expected {
						t.Errorf("Expected %s for %s/%s, got %s", tt.expected, tt.goos, tt.goarch, platform)
					}
				}
			}
		})
	}
}

// TestReleaseAssetNameFormat tests the release asset name format more thoroughly
func TestReleaseAssetNameFormat(t *testing.T) {
	platformUtil := PlatformUtil{}
	
	testVersions := []string{
		"1.0.0",
		"2.1.3",
		"0.9.15",
		"1.0.0-alpha",
		"1.0.0-beta.1",
		"1.0.0-rc.2",
		"1.2.3-dev+build.456",
	}

	for _, version := range testVersions {
		t.Run("version_"+version, func(t *testing.T) {
			assetName := platformUtil.GetReleaseAssetName(version)
			
			// Should follow pattern: lpm-{version}-{platform}.zip
			parts := strings.Split(assetName, "-")
			if len(parts) < 3 {
				t.Errorf("Asset name %s should have at least 3 parts separated by '-'", assetName)
				return
			}
			
			if parts[0] != "lpm" {
				t.Errorf("Asset name %s should start with 'lpm'", assetName)
			}
			
			if !strings.HasSuffix(assetName, ".zip") {
				t.Errorf("Asset name %s should end with '.zip'", assetName)
			}
			
			// Version should be in the name
			if !strings.Contains(assetName, version) {
				t.Errorf("Asset name %s should contain version %s", assetName, version)
			}
		})
	}
}