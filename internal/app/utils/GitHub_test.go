package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/v39/github"
)

func TestNewGitHubUtil(t *testing.T) {
	// Test creating GitHubUtil without authentication
	originalToken := os.Getenv("GITHUB_PAT")
	os.Unsetenv("GITHUB_PAT")
	defer func() {
		if originalToken != "" {
			os.Setenv("GITHUB_PAT", originalToken)
		}
	}()

	githubUtil := NewGitHubUtil()
	if githubUtil == nil {
		t.Error("NewGitHubUtil() should not return nil")
	}
	if githubUtil.client == nil {
		t.Error("GitHubUtil client should not be nil")
	}
	if githubUtil.ctx == nil {
		t.Error("GitHubUtil context should not be nil")
	}

	// Test creating GitHubUtil with authentication
	os.Setenv("GITHUB_PAT", "test-token")
	githubUtilAuth := NewGitHubUtil()
	if githubUtilAuth == nil {
		t.Error("NewGitHubUtil() with auth should not return nil")
	}
	if githubUtilAuth.client == nil {
		t.Error("GitHubUtil client with auth should not be nil")
	}
}

func TestGetLatestRelease(t *testing.T) {
	githubUtil := NewGitHubUtil()

	// Test with valid public repository (using a well-known repo)
	t.Run("valid repository", func(t *testing.T) {
		// Use a repository that's likely to have releases
		release, err := githubUtil.GetLatestRelease("golang", "go")
		if err != nil {
			// This test may fail due to rate limiting in CI, so we'll make it less strict
			t.Logf("Warning: Could not get latest release for golang/go: %v", err)
			if !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
				t.Errorf("Unexpected error: %v", err)
			}
			return
		}
		if release == nil {
			t.Error("Release should not be nil")
		}
		if release.GetTagName() == "" {
			t.Error("Release tag name should not be empty")
		}
	})

	// Test with empty owner
	t.Run("empty owner", func(t *testing.T) {
		_, err := githubUtil.GetLatestRelease("", "repo")
		if err == nil {
			t.Error("Expected error for empty owner")
		}
		if !strings.Contains(err.Error(), "owner and repo cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with empty repo
	t.Run("empty repo", func(t *testing.T) {
		_, err := githubUtil.GetLatestRelease("owner", "")
		if err == nil {
			t.Error("Expected error for empty repo")
		}
		if !strings.Contains(err.Error(), "owner and repo cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with non-existent repository
	t.Run("non-existent repository", func(t *testing.T) {
		_, err := githubUtil.GetLatestRelease("nonexistent-owner-12345", "nonexistent-repo-12345")
		if err == nil {
			t.Error("Expected error for non-existent repository")
		}
		// Error message should indicate the repository was not found
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
			t.Errorf("Expected 'not found' error or rate limit, got: %v", err)
		}
	})
}

func TestGetRelease(t *testing.T) {
	githubUtil := NewGitHubUtil()

	// Test with valid parameters
	t.Run("valid parameters", func(t *testing.T) {
		// Use a known release version - we'll use a Go release
		release, err := githubUtil.GetRelease("golang", "go", "go1.20.0")
		if err != nil {
			// This test may fail due to rate limiting in CI
			t.Logf("Warning: Could not get release for golang/go: %v", err)
			if !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
				t.Errorf("Unexpected error: %v", err)
			}
			return
		}
		if release == nil {
			t.Error("Release should not be nil")
		}
		if release.GetTagName() == "" {
			t.Error("Release tag name should not be empty")
		}
	})

	// Test with empty owner
	t.Run("empty owner", func(t *testing.T) {
		_, err := githubUtil.GetRelease("", "repo", "v1.0.0")
		if err == nil {
			t.Error("Expected error for empty owner")
		}
		if !strings.Contains(err.Error(), "owner, repo, and version cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with empty repo
	t.Run("empty repo", func(t *testing.T) {
		_, err := githubUtil.GetRelease("owner", "", "v1.0.0")
		if err == nil {
			t.Error("Expected error for empty repo")
		}
		if !strings.Contains(err.Error(), "owner, repo, and version cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with empty version
	t.Run("empty version", func(t *testing.T) {
		_, err := githubUtil.GetRelease("owner", "repo", "")
		if err == nil {
			t.Error("Expected error for empty version")
		}
		if !strings.Contains(err.Error(), "owner, repo, and version cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with non-existent version
	t.Run("non-existent version", func(t *testing.T) {
		_, err := githubUtil.GetRelease("golang", "go", "v999.999.999")
		if err == nil {
			t.Error("Expected error for non-existent version")
		}
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
			t.Errorf("Expected 'not found' error or rate limit, got: %v", err)
		}
	})
}

func TestGetReleaseAsset(t *testing.T) {
	// Create a mock release with assets
	mockRelease := &github.RepositoryRelease{
		Assets: []*github.ReleaseAsset{
			{
				Name: github.String("lpm-1.0.0-linux.zip"),
			},
			{
				Name: github.String("lpm-1.0.0-darwin.zip"),
			},
			{
				Name: github.String("lpm-1.0.0-windows.zip"),
			},
		},
	}

	githubUtil := NewGitHubUtil()

	// Test finding exact asset name
	t.Run("exact asset name", func(t *testing.T) {
		asset, err := githubUtil.GetReleaseAsset(mockRelease, "lpm-1.0.0-linux.zip")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if asset == nil {
			t.Error("Asset should not be nil")
		}
		if asset.GetName() != "lpm-1.0.0-linux.zip" {
			t.Errorf("Expected asset name 'lpm-1.0.0-linux.zip', got '%s'", asset.GetName())
		}
	})

	// Test finding asset with partial name
	t.Run("partial asset name", func(t *testing.T) {
		asset, err := githubUtil.GetReleaseAsset(mockRelease, "darwin")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if asset == nil {
			t.Error("Asset should not be nil")
		}
		if !strings.Contains(asset.GetName(), "darwin") {
			t.Errorf("Expected asset name to contain 'darwin', got '%s'", asset.GetName())
		}
	})

	// Test with nil release
	t.Run("nil release", func(t *testing.T) {
		_, err := githubUtil.GetReleaseAsset(nil, "test-asset")
		if err == nil {
			t.Error("Expected error for nil release")
		}
		if !strings.Contains(err.Error(), "release cannot be nil") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with empty asset name
	t.Run("empty asset name", func(t *testing.T) {
		_, err := githubUtil.GetReleaseAsset(mockRelease, "")
		if err == nil {
			t.Error("Expected error for empty asset name")
		}
		if !strings.Contains(err.Error(), "asset name cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with non-existent asset
	t.Run("non-existent asset", func(t *testing.T) {
		_, err := githubUtil.GetReleaseAsset(mockRelease, "non-existent-asset")
		if err == nil {
			t.Error("Expected error for non-existent asset")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
		// Should include list of available assets
		if !strings.Contains(err.Error(), "Available assets:") {
			t.Errorf("Error should include available assets list, got: %v", err)
		}
	})
}

func TestDownloadReleaseAsset(t *testing.T) {
	githubUtil := NewGitHubUtil()

	// Test with nil asset
	t.Run("nil asset", func(t *testing.T) {
		_, err := githubUtil.DownloadReleaseAsset(nil)
		if err == nil {
			t.Error("Expected error for nil asset")
		}
		if !strings.Contains(err.Error(), "asset cannot be nil") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with asset without download URL
	t.Run("asset without download URL", func(t *testing.T) {
		mockAsset := &github.ReleaseAsset{
			Name: github.String("test-asset"),
			// No BrowserDownloadURL set
		}
		_, err := githubUtil.DownloadReleaseAsset(mockAsset)
		if err == nil {
			t.Error("Expected error for asset without download URL")
		}
		if !strings.Contains(err.Error(), "does not have a download URL") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with invalid URL
	t.Run("invalid download URL", func(t *testing.T) {
		mockAsset := &github.ReleaseAsset{
			Name:               github.String("test-asset"),
			BrowserDownloadURL: github.String("invalid-url"),
		}
		_, err := githubUtil.DownloadReleaseAsset(mockAsset)
		if err == nil {
			t.Error("Expected error for invalid download URL")
		}
		// Should contain some form of network or URL error
		if !strings.Contains(err.Error(), "failed to") {
			t.Errorf("Expected failure error, got: %v", err)
		}
	})

	// Note: We can't easily test successful downloads without mocking HTTP client
	// or using a real asset, which would make tests dependent on external resources
}

func TestListReleases(t *testing.T) {
	githubUtil := NewGitHubUtil()

	// Test with valid repository
	t.Run("valid repository", func(t *testing.T) {
		releases, err := githubUtil.ListReleases("golang", "go")
		if err != nil {
			// This test may fail due to rate limiting in CI
			t.Logf("Warning: Could not list releases for golang/go: %v", err)
			if !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
				t.Errorf("Unexpected error: %v", err)
			}
			return
		}
		if releases == nil {
			t.Error("Releases should not be nil")
		}
		// golang/go should have many releases
		if len(releases) == 0 {
			t.Error("Expected at least one release for golang/go")
		}
	})

	// Test with empty owner
	t.Run("empty owner", func(t *testing.T) {
		_, err := githubUtil.ListReleases("", "repo")
		if err == nil {
			t.Error("Expected error for empty owner")
		}
		if !strings.Contains(err.Error(), "owner and repo cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with empty repo
	t.Run("empty repo", func(t *testing.T) {
		_, err := githubUtil.ListReleases("owner", "")
		if err == nil {
			t.Error("Expected error for empty repo")
		}
		if !strings.Contains(err.Error(), "owner and repo cannot be empty") {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	// Test with non-existent repository
	t.Run("non-existent repository", func(t *testing.T) {
		_, err := githubUtil.ListReleases("nonexistent-owner-12345", "nonexistent-repo-12345")
		if err == nil {
			t.Error("Expected error for non-existent repository")
		}
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
			t.Errorf("Expected 'not found' error or rate limit, got: %v", err)
		}
	})
}

func TestGetRateLimitInfo(t *testing.T) {
	githubUtil := NewGitHubUtil()

	// Test getting rate limit info
	t.Run("get rate limit info", func(t *testing.T) {
		limits, err := githubUtil.GetRateLimitInfo()
		if err != nil {
			// This test may fail due to network issues or authentication
			t.Logf("Warning: Could not get rate limit info: %v", err)
			if !strings.Contains(err.Error(), "rate limit") && !strings.Contains(err.Error(), "403") {
				t.Errorf("Unexpected error: %v", err)
			}
			return
		}
		if limits == nil {
			t.Error("Rate limits should not be nil")
		}
		// Core rate limit should exist
		if limits.Core == nil {
			t.Error("Core rate limit should not be nil")
		}
	})
}

// Integration test helper - only runs if GITHUB_INTEGRATION_TEST env var is set
func TestGitHubUtilIntegration(t *testing.T) {
	if os.Getenv("GITHUB_INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set GITHUB_INTEGRATION_TEST=1 to run.")
	}

	githubUtil := NewGitHubUtil()

	// Test full workflow with a known repository
	t.Run("full workflow", func(t *testing.T) {
		// Get latest release
		release, err := githubUtil.GetLatestRelease("golang", "go")
		if err != nil {
			t.Fatalf("Failed to get latest release: %v", err)
		}

		// List all releases
		releases, err := githubUtil.ListReleases("golang", "go")
		if err != nil {
			t.Fatalf("Failed to list releases: %v", err)
		}

		if len(releases) == 0 {
			t.Error("Expected at least one release")
		}

		// Check if latest release is in the list
		found := false
		for _, r := range releases {
			if r.GetTagName() == release.GetTagName() {
				found = true
				break
			}
		}

		if !found {
			t.Error("Latest release should be in the releases list")
		}

		// Get specific release
		specificRelease, err := githubUtil.GetRelease("golang", "go", release.GetTagName())
		if err != nil {
			t.Fatalf("Failed to get specific release: %v", err)
		}

		if specificRelease.GetTagName() != release.GetTagName() {
			t.Errorf("Expected same tag name, got %s vs %s", specificRelease.GetTagName(), release.GetTagName())
		}

		fmt.Printf("Integration test successful. Latest Go release: %s\n", release.GetTagName())
	})
}