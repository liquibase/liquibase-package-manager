package utils

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		versionStr  string
		wantErr     bool
		wantVersion string
	}{
		{
			name:        "simple version",
			versionStr:  "0.2.9",
			wantErr:     false,
			wantVersion: "0.2.9",
		},
		{
			name:        "version with v prefix",
			versionStr:  "v0.2.9",
			wantErr:     false,
			wantVersion: "0.2.9",
		},
		{
			name:        "version with whitespace",
			versionStr:  "  0.2.9  ",
			wantErr:     false,
			wantVersion: "0.2.9",
		},
		{
			name:        "pre-release version",
			versionStr:  "1.0.0-alpha",
			wantErr:     false,
			wantVersion: "1.0.0-alpha",
		},
		{
			name:        "pre-release with v prefix",
			versionStr:  "v1.0.0-beta.1",
			wantErr:     false,
			wantVersion: "1.0.0-beta.1",
		},
		{
			name:        "version with metadata",
			versionStr:  "1.0.0+build.123",
			wantErr:     false,
			wantVersion: "1.0.0+build.123",
		},
		{
			name:       "invalid version",
			versionStr: "invalid",
			wantErr:    true,
		},
		{
			name:       "empty version",
			versionStr: "",
			wantErr:    true,
		},
		{
			name:        "major.minor only",
			versionStr:  "1.0",
			wantErr:     false,
			wantVersion: "1.0.0",
		},
		{
			name:        "major only",
			versionStr:  "1",
			wantErr:     false,
			wantVersion: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.versionStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.wantVersion {
				t.Errorf("ParseVersion() = %v, want %v", got.String(), tt.wantVersion)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    int
		wantErr bool
	}{
		{
			name:    "current less than latest",
			current: "0.2.8",
			latest:  "0.2.9",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "current equals latest",
			current: "0.2.9",
			latest:  "0.2.9",
			want:    0,
			wantErr: false,
		},
		{
			name:    "current greater than latest",
			current: "0.3.0",
			latest:  "0.2.9",
			want:    1,
			wantErr: false,
		},
		{
			name:    "major version difference",
			current: "1.0.0",
			latest:  "2.0.0",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "minor version difference",
			current: "1.1.0",
			latest:  "1.2.0",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "patch version difference",
			current: "1.0.1",
			latest:  "1.0.2",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "with v prefix on both",
			current: "v0.2.8",
			latest:  "v0.2.9",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "mixed v prefix",
			current: "0.2.8",
			latest:  "v0.2.9",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "pre-release versions",
			current: "1.0.0-alpha",
			latest:  "1.0.0-beta",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "pre-release vs stable",
			current: "1.0.0-rc.1",
			latest:  "1.0.0",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "invalid current version",
			current: "invalid",
			latest:  "1.0.0",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid latest version",
			current: "1.0.0",
			latest:  "invalid",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareVersions(tt.current, tt.latest)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUpdateAvailable(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    bool
		wantErr bool
	}{
		{
			name:    "update available",
			current: "0.2.8",
			latest:  "0.2.9",
			want:    true,
			wantErr: false,
		},
		{
			name:    "no update - same version",
			current: "0.2.9",
			latest:  "0.2.9",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no update - current newer",
			current: "0.3.0",
			latest:  "0.2.9",
			want:    false,
			wantErr: false,
		},
		{
			name:    "major update available",
			current: "1.5.3",
			latest:  "2.0.0",
			want:    true,
			wantErr: false,
		},
		{
			name:    "minor update available",
			current: "1.2.3",
			latest:  "1.3.0",
			want:    true,
			wantErr: false,
		},
		{
			name:    "patch update available",
			current: "1.2.3",
			latest:  "1.2.4",
			want:    true,
			wantErr: false,
		},
		{
			name:    "pre-release to stable update",
			current: "1.0.0-rc.1",
			latest:  "1.0.0",
			want:    true,
			wantErr: false,
		},
		{
			name:    "pre-release to pre-release update",
			current: "1.0.0-alpha",
			latest:  "1.0.0-beta",
			want:    true,
			wantErr: false,
		},
		{
			name:    "no update - stable to pre-release",
			current: "1.0.0",
			latest:  "1.0.0-rc.1",
			want:    false,
			wantErr: false,
		},
		{
			name:    "with v prefix",
			current: "v0.2.8",
			latest:  "v0.2.9",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid current version",
			current: "invalid",
			latest:  "1.0.0",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid latest version",
			current: "1.0.0",
			latest:  "invalid",
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty versions",
			current: "",
			latest:  "",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsUpdateAvailable(tt.current, tt.latest)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUpdateAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUpdateAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentVersion(t *testing.T) {
	// GetCurrentVersion returns the embedded version from app.Version()
	// We can't easily mock this in a unit test, but we can verify it returns
	// a non-empty string that can be parsed as a version
	version := GetCurrentVersion()
	
	if version == "" {
		t.Error("GetCurrentVersion() returned empty string")
	}
	
	// Verify the returned version can be parsed
	_, err := ParseVersion(version)
	if err != nil {
		t.Errorf("GetCurrentVersion() returned unparseable version: %v", err)
	}
}