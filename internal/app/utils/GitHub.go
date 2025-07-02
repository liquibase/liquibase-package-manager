package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// GitHubUtil provides GitHub API functionality for release management
type GitHubUtil struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubUtil creates a new GitHubUtil instance with optional authentication
func NewGitHubUtil() *GitHubUtil {
	ctx := context.Background()
	var client *github.Client

	// Check for GitHub PAT for authentication
	if token := os.Getenv("GITHUB_PAT"); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		// Use unauthenticated client (with rate limits)
		client = github.NewClient(nil)
	}

	return &GitHubUtil{
		client: client,
		ctx:    ctx,
	}
}

// GetLatestRelease retrieves the latest release for a GitHub repository
func (g *GitHubUtil) GetLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo cannot be empty")
	}

	fmt.Printf("[DEBUG] GitHub API: Requesting latest release for %s/%s\n", owner, repo)
	release, resp, err := g.client.Repositories.GetLatestRelease(g.ctx, owner, repo)
	
	// Log rate limit information if available
	if resp != nil {
		fmt.Printf("[DEBUG] GitHub API: Rate limit remaining: %d/%d\n", resp.Rate.Remaining, resp.Rate.Limit)
	}
	
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			fmt.Printf("[DEBUG] GitHub API: Repository %s/%s has no releases (404)\n", owner, repo)
			return nil, fmt.Errorf("no releases found for repository %s/%s", owner, repo)
		}
		if resp != nil && resp.StatusCode == 403 {
			fmt.Printf("[DEBUG] GitHub API: Rate limited or access denied for %s/%s (403)\n", owner, repo)
			return nil, fmt.Errorf("rate limit exceeded or access denied for repository %s/%s", owner, repo)
		}
		fmt.Printf("[DEBUG] GitHub API: Request failed for %s/%s: %v\n", owner, repo, err)
		return nil, fmt.Errorf("failed to get latest release for %s/%s: %w", owner, repo, err)
	}

	if release == nil {
		fmt.Printf("[DEBUG] GitHub API: No release data returned for %s/%s\n", owner, repo)
		return nil, fmt.Errorf("no latest release found for repository %s/%s", owner, repo)
	}

	fmt.Printf("[DEBUG] GitHub API: Found latest release %s for %s/%s\n", release.GetTagName(), owner, repo)
	return release, nil
}

// GetRelease retrieves a specific release by version/tag for a GitHub repository
func (g *GitHubUtil) GetRelease(owner, repo, version string) (*github.RepositoryRelease, error) {
	if owner == "" || repo == "" || version == "" {
		return nil, fmt.Errorf("owner, repo, and version cannot be empty")
	}

	// Ensure version starts with 'v' prefix if not already present for tag lookup
	tag := version
	if !strings.HasPrefix(version, "v") {
		tag = "v" + version
	}

	fmt.Printf("[DEBUG] GitHub API: Requesting release %s (tag: %s) for %s/%s\n", version, tag, owner, repo)
	release, resp, err := g.client.Repositories.GetReleaseByTag(g.ctx, owner, repo, tag)
	
	// Log rate limit information if available
	if resp != nil {
		fmt.Printf("[DEBUG] GitHub API: Rate limit remaining: %d/%d\n", resp.Rate.Remaining, resp.Rate.Limit)
	}
	
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			fmt.Printf("[DEBUG] GitHub API: Release %s not found, trying without 'v' prefix\n", tag)
			// Try without 'v' prefix if the first attempt failed
			if strings.HasPrefix(tag, "v") {
				release, resp, err = g.client.Repositories.GetReleaseByTag(g.ctx, owner, repo, version)
				if err != nil {
					if resp != nil && resp.StatusCode == 404 {
						fmt.Printf("[DEBUG] GitHub API: Release %s not found for %s/%s (tried both %s and %s)\n", version, owner, repo, tag, version)
						return nil, fmt.Errorf("release %s not found for repository %s/%s", version, owner, repo)
					}
				} else {
					fmt.Printf("[DEBUG] GitHub API: Found release %s using tag %s\n", version, version)
					return release, nil
				}
			} else {
				fmt.Printf("[DEBUG] GitHub API: Release %s not found for %s/%s\n", version, owner, repo)
				return nil, fmt.Errorf("release %s not found for repository %s/%s", version, owner, repo)
			}
		}
		if resp != nil && resp.StatusCode == 403 {
			fmt.Printf("[DEBUG] GitHub API: Rate limited or access denied for %s/%s (403)\n", owner, repo)
			return nil, fmt.Errorf("rate limit exceeded or access denied for repository %s/%s", owner, repo)
		}
		if err != nil {
			fmt.Printf("[DEBUG] GitHub API: Request failed for %s/%s (version %s): %v\n", owner, repo, version, err)
			return nil, fmt.Errorf("failed to get release %s for %s/%s: %w", version, owner, repo, err)
		}
	}

	if release == nil {
		fmt.Printf("[DEBUG] GitHub API: No release data returned for %s/%s (version %s)\n", owner, repo, version)
		return nil, fmt.Errorf("release %s not found for repository %s/%s", version, owner, repo)
	}

	fmt.Printf("[DEBUG] GitHub API: Found release %s for %s/%s\n", release.GetTagName(), owner, repo)
	return release, nil
}

// GetReleaseAsset finds a specific asset by name within a release
func (g *GitHubUtil) GetReleaseAsset(release *github.RepositoryRelease, assetName string) (*github.ReleaseAsset, error) {
	if release == nil {
		return nil, fmt.Errorf("release cannot be nil")
	}
	if assetName == "" {
		return nil, fmt.Errorf("asset name cannot be empty")
	}

	for _, asset := range release.Assets {
		if asset.GetName() == assetName {
			return asset, nil
		}
	}

	// If exact match not found, try partial matching for common patterns
	for _, asset := range release.Assets {
		if strings.Contains(asset.GetName(), assetName) {
			return asset, nil
		}
	}

	availableAssets := make([]string, len(release.Assets))
	for i, asset := range release.Assets {
		availableAssets[i] = asset.GetName()
	}

	return nil, fmt.Errorf("asset '%s' not found in release. Available assets: %s", assetName, strings.Join(availableAssets, ", "))
}

// DownloadReleaseAsset downloads the content of a release asset
func (g *GitHubUtil) DownloadReleaseAsset(asset *github.ReleaseAsset) ([]byte, error) {
	if asset == nil {
		return nil, fmt.Errorf("asset cannot be nil")
	}

	downloadURL := asset.GetBrowserDownloadURL()
	if downloadURL == "" {
		return nil, fmt.Errorf("asset does not have a download URL")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Minute, // 5 minute timeout for large downloads
	}

	req, err := http.NewRequestWithContext(g.ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	// Add GitHub token to request if available for better rate limits
	if token := os.Getenv("GITHUB_PAT"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset %s: %w", asset.GetName(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download asset %s: HTTP %d", asset.GetName(), resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset %s content: %w", asset.GetName(), err)
	}

	return data, nil
}

// ListReleases retrieves all releases for a repository (useful for version listing)
func (g *GitHubUtil) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo cannot be empty")
	}

	var allReleases []*github.RepositoryRelease
	opts := &github.ListOptions{
		Page:    1,
		PerPage: 100, // Maximum per page
	}

	for {
		releases, resp, err := g.client.Repositories.ListReleases(g.ctx, owner, repo, opts)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return nil, fmt.Errorf("repository %s/%s not found", owner, repo)
			}
			if resp != nil && resp.StatusCode == 403 {
				return nil, fmt.Errorf("rate limit exceeded or access denied for repository %s/%s", owner, repo)
			}
			return nil, fmt.Errorf("failed to list releases for %s/%s: %w", owner, repo, err)
		}

		allReleases = append(allReleases, releases...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReleases, nil
}

// GetRateLimitInfo returns the current rate limit status
func (g *GitHubUtil) GetRateLimitInfo() (*github.RateLimits, error) {
	limits, _, err := g.client.RateLimits(g.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limits: %w", err)
	}
	return limits, nil
}