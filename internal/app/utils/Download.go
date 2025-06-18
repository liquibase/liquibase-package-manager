package utils

import (
	"archive/zip"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
)

// DownloadUtil provides secure download functionality with progress tracking and checksum verification
type DownloadUtil struct {
	client *http.Client
}

// ProgressWriter wraps an io.Writer to provide download progress feedback
type ProgressWriter struct {
	writer     io.Writer
	total      int64
	written    int64
	onProgress func(current, total int64)
}

// NewProgressWriter creates a new progress writer with callback
func NewProgressWriter(writer io.Writer, total int64, onProgress func(current, total int64)) *ProgressWriter {
	return &ProgressWriter{
		writer:     writer,
		total:      total,
		onProgress: onProgress,
	}
}

// Write implements io.Writer interface with progress tracking
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if n > 0 {
		pw.written += int64(n)
		if pw.onProgress != nil {
			pw.onProgress(pw.written, pw.total)
		}
	}
	return n, err
}

// NewDownloadUtil creates a new DownloadUtil with secure HTTP client configuration
func NewDownloadUtil() *DownloadUtil {
	// Create secure HTTP client with proper TLS configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false, // Always verify certificates
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Minute, // Generous timeout for large binary downloads
	}

	return &DownloadUtil{
		client: client,
	}
}

// DownloadWithProgress downloads a file from URL to filepath with progress indication
func (d *DownloadUtil) DownloadWithProgress(downloadURL, filePath string) error {
	if downloadURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	if filePath == "" {
		return fmt.Errorf("filepath cannot be empty")
	}

	// Validate URL format
	if _, err := url.Parse(downloadURL); err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Create temporary file for atomic operation
	tempFile := filePath + ".tmp"
	defer os.Remove(tempFile) // Clean up on any failure

	// Create the request
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add User-Agent to identify our client
	req.Header.Set("User-Agent", "liquibase-package-manager/1.0")

	// Add GitHub token if downloading from GitHub and token is available
	if strings.Contains(downloadURL, "github.com") || strings.Contains(downloadURL, "githubusercontent.com") {
		if token := os.Getenv("GITHUB_PAT"); token != "" {
			req.Header.Set("Authorization", "token "+token)
		}
	}

	// Execute the request
	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	// Create output file
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Get content length for progress tracking
	contentLength := resp.ContentLength
	if contentLength > 0 {
		fmt.Printf("Downloading %s (%s)...\n", filepath.Base(filePath), formatBytes(contentLength))
	} else {
		fmt.Printf("Downloading %s...\n", filepath.Base(filePath))
	}

	// Create progress writer
	var lastPercent int
	progressWriter := NewProgressWriter(out, contentLength, func(current, total int64) {
		if total > 0 {
			percent := int((current * 100) / total)
			if percent != lastPercent && percent%10 == 0 {
				fmt.Printf("Progress: %d%% (%s/%s)\n", percent, formatBytes(current), formatBytes(total))
				lastPercent = percent
			}
		}
	})

	// Copy with progress tracking
	_, err = io.Copy(progressWriter, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Atomic move from temp to final location
	if err := os.Rename(tempFile, filePath); err != nil {
		return fmt.Errorf("failed to move downloaded file to final location: %w", err)
	}

	fmt.Printf("Download completed: %s\n", filePath)
	return nil
}

// VerifyChecksum verifies the checksum of a file using the specified algorithm
func (d *DownloadUtil) VerifyChecksum(filePath, expectedChecksum, algorithm string) error {
	if filePath == "" {
		return fmt.Errorf("filepath cannot be empty")
	}
	if expectedChecksum == "" {
		return fmt.Errorf("expected checksum cannot be empty")
	}
	if algorithm == "" {
		return fmt.Errorf("algorithm cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum verification: %w", err)
	}
	defer file.Close()

	// Create appropriate hasher
	var hasher hash.Hash
	switch strings.ToLower(algorithm) {
	case "sha256":
		hasher = sha256.New()
	case "sha1":
		hasher = sha1.New()
	default:
		return fmt.Errorf("unsupported checksum algorithm: %s (supported: sha256, sha1)", algorithm)
	}

	// Calculate checksum
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Compare checksums
	actualChecksum := fmt.Sprintf("%x", hasher.Sum(nil))
	expectedChecksum = strings.ToLower(strings.TrimSpace(expectedChecksum))

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum verification failed: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	fmt.Printf("Checksum verified (%s): %s\n", algorithm, actualChecksum)
	return nil
}

// DownloadToTemp downloads a file to the system temporary directory
func (d *DownloadUtil) DownloadToTemp(downloadURL, filename string) (string, error) {
	if downloadURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	// Create temporary directory for our downloads
	tempDir := filepath.Join(os.TempDir(), "lpm-downloads")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Generate unique temporary file path
	tempFile := filepath.Join(tempDir, filename)

	// Download to temporary location
	if err := d.DownloadWithProgress(downloadURL, tempFile); err != nil {
		return "", fmt.Errorf("failed to download to temp directory: %w", err)
	}

	return tempFile, nil
}

// DownloadBinaryUpdate downloads an lpm binary update from a GitHub release
func (d *DownloadUtil) DownloadBinaryUpdate(release *github.RepositoryRelease, platform string) (string, error) {
	if release == nil {
		return "", fmt.Errorf("release cannot be nil")
	}
	if platform == "" {
		platform = runtime.GOOS
	}

	fmt.Printf("[DEBUG] Download: Searching for binary asset for platform %s in release %s\n", platform, release.GetTagName())

	// Determine asset name pattern based on platform
	var assetPattern string
	switch platform {
	case "darwin":
		assetPattern = "darwin"
	case "linux":
		assetPattern = "linux"
	case "windows":
		assetPattern = "windows"
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	fmt.Printf("[DEBUG] Download: Looking for assets containing pattern '%s'\n", assetPattern)

	// Find the appropriate asset
	var targetAsset *github.ReleaseAsset
	for _, asset := range release.Assets {
		assetName := strings.ToLower(asset.GetName())
		fmt.Printf("[DEBUG] Download: Checking asset %s\n", asset.GetName())
		if strings.Contains(assetName, assetPattern) {
			// Prefer exact platform matches and common binary formats
			if strings.HasSuffix(assetName, ".zip") || strings.HasSuffix(assetName, ".tar.gz") {
				fmt.Printf("[DEBUG] Download: Selected asset %s for platform %s\n", asset.GetName(), platform)
				targetAsset = asset
				break
			}
		}
	}

	if targetAsset == nil {
		// List available assets for debugging
		var availableAssets []string
		for _, asset := range release.Assets {
			availableAssets = append(availableAssets, asset.GetName())
		}
		fmt.Printf("[DEBUG] Download: No suitable asset found. Available assets: %s\n", strings.Join(availableAssets, ", "))
		return "", fmt.Errorf("no suitable binary asset found for platform %s. Available assets: %s", platform, strings.Join(availableAssets, ", "))
	}

	// Download the asset
	downloadURL := targetAsset.GetBrowserDownloadURL()
	if downloadURL == "" {
		return "", fmt.Errorf("asset does not have a download URL")
	}

	filename := targetAsset.GetName()
	fmt.Printf("[DEBUG] Download: Starting download of %s from %s\n", filename, downloadURL)
	fmt.Printf("[DEBUG] Download: Asset size: %s\n", formatBytes(int64(targetAsset.GetSize())))
	
	tempFile, err := d.DownloadToTemp(downloadURL, filename)
	if err != nil {
		fmt.Printf("[ERROR] Download: Failed to download %s: %v\n", filename, err)
		return "", fmt.Errorf("failed to download binary update: %w", err)
	}
	fmt.Printf("[DEBUG] Download: Successfully downloaded %s to %s\n", filename, tempFile)

	// If it's a ZIP file, extract the binary
	if strings.HasSuffix(strings.ToLower(filename), ".zip") {
		fmt.Printf("[DEBUG] Download: Extracting binary from ZIP archive %s\n", filename)
		extractedPath, err := d.extractBinaryFromZip(tempFile, platform)
		if err != nil {
			fmt.Printf("[ERROR] Download: Failed to extract binary from ZIP: %v\n", err)
			os.Remove(tempFile) // Clean up downloaded ZIP
			return "", fmt.Errorf("failed to extract binary from ZIP: %w", err)
		}
		fmt.Printf("[DEBUG] Download: Binary extracted to %s\n", extractedPath)
		os.Remove(tempFile) // Clean up ZIP file after extraction
		return extractedPath, nil
	}

	return tempFile, nil
}

// extractBinaryFromZip extracts the main binary from a ZIP archive
func (d *DownloadUtil) extractBinaryFromZip(zipPath, platform string) (string, error) {
	fmt.Printf("[DEBUG] Extract: Opening ZIP file %s\n", zipPath)
	// Open ZIP file
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer r.Close()

	// Determine expected binary name
	expectedBinary := "lpm"
	if platform == "windows" {
		expectedBinary = "lpm.exe"
	}
	fmt.Printf("[DEBUG] Extract: Looking for binary file '%s' in ZIP\n", expectedBinary)

	// Find the binary file in the ZIP
	var binaryFile *zip.File
	for _, f := range r.File {
		fileName := filepath.Base(f.Name)
		fmt.Printf("[DEBUG] Extract: Found file in ZIP: %s\n", f.Name)
		if fileName == expectedBinary {
			fmt.Printf("[DEBUG] Extract: Exact match found: %s\n", f.Name)
			binaryFile = f
			break
		}
		// Also check for common patterns like "lpm-*" or in subdirectories
		if strings.Contains(fileName, "lpm") && (strings.HasSuffix(fileName, ".exe") || !strings.Contains(fileName, ".")) {
			fmt.Printf("[DEBUG] Extract: Potential match found: %s\n", f.Name)
			binaryFile = f
		}
	}

	if binaryFile == nil {
		// List ZIP contents for debugging
		var contents []string
		for _, f := range r.File {
			contents = append(contents, f.Name)
		}
		fmt.Printf("[DEBUG] Extract: Binary %s not found. ZIP contents: %s\n", expectedBinary, strings.Join(contents, ", "))
		return "", fmt.Errorf("binary file %s not found in ZIP. Contents: %s", expectedBinary, strings.Join(contents, ", "))
	}

	// Create extraction path
	extractDir := filepath.Dir(zipPath)
	extractedPath := filepath.Join(extractDir, expectedBinary)

	// Open binary file in ZIP
	rc, err := binaryFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open binary file in ZIP: %w", err)
	}
	defer rc.Close()

	// Create output file
	outFile, err := os.Create(extractedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create extracted binary file: %w", err)
	}
	defer outFile.Close()

	// Copy binary content
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return "", fmt.Errorf("failed to extract binary: %w", err)
	}

	// Make binary executable on Unix-like systems
	if platform != "windows" {
		if err := os.Chmod(extractedPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	fmt.Printf("Binary extracted: %s\n", extractedPath)
	return extractedPath, nil
}

// formatBytes formats byte counts as human readable strings
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// DownloadWithChecksum downloads a file and verifies its checksum in one operation
func (d *DownloadUtil) DownloadWithChecksum(downloadURL, filePath, expectedChecksum, algorithm string) error {
	// Download the file
	if err := d.DownloadWithProgress(downloadURL, filePath); err != nil {
		return err
	}

	// Verify checksum if provided
	if expectedChecksum != "" && algorithm != "" {
		if err := d.VerifyChecksum(filePath, expectedChecksum, algorithm); err != nil {
			// Clean up downloaded file on checksum failure
			os.Remove(filePath)
			return err
		}
	}

	return nil
}

// CleanupTempDownloads removes temporary download files older than specified duration
func CleanupTempDownloads(maxAge time.Duration) error {
	tempDir := filepath.Join(os.TempDir(), "lpm-downloads")
	
	// Check if temp directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return nil // Nothing to clean up
	}

	// Read directory contents
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	now := time.Now()
	cleaned := 0

	for _, entry := range entries {
		filePath := filepath.Join(tempDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue // Skip files we can't stat
		}

		// Remove files older than maxAge
		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(filePath); err == nil {
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		fmt.Printf("Cleaned up %d temporary download files\n", cleaned)
	}

	return nil
}